package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pfarrer/foghorn/logger"
)

type CheckConfig interface {
	GetName() string
	GetSchedule() string
	IsEnabled() bool
}

type IntervalCheckConfig interface {
	CheckConfig
	GetScheduleType() ScheduleType
	GetInterval() string
}

type CheckExecutor interface {
	Execute(check CheckConfig) error
}

type ScheduledCheck struct {
	Config       CheckConfig
	NextRun      time.Time
	LastRun      *time.Time
	Running      bool
	ScheduleType ScheduleType
	Interval     time.Duration
}

type Scheduler struct {
	checks              map[string]*ScheduledCheck
	executor            CheckExecutor
	ticker              *time.Ticker
	stopChan            chan struct{}
	location            *time.Location
	maxConcurrentChecks int
	runningChecks       int
	queue               []CheckConfig
}

func NewScheduler(executor CheckExecutor, location *time.Location, maxConcurrentChecks int) *Scheduler {
	if location == nil {
		location = time.UTC
	}
	return &Scheduler{
		checks:              make(map[string]*ScheduledCheck),
		executor:            executor,
		stopChan:            make(chan struct{}),
		location:            location,
		maxConcurrentChecks: maxConcurrentChecks,
		queue:               make([]CheckConfig, 0),
	}
}

func (s *Scheduler) AddCheck(config CheckConfig) error {
	if config.GetSchedule() == "" {
		return fmt.Errorf("check %s: schedule is required", config.GetName())
	}

	var nextRun time.Time
	var scheduleType ScheduleType
	var interval time.Duration
	var err error

	if intervalCheck, ok := config.(IntervalCheckConfig); ok && intervalCheck.GetScheduleType() == ScheduleTypeInterval {
		scheduleType = ScheduleTypeInterval
		interval, err = parseInterval(intervalCheck.GetInterval())
		if err != nil {
			return fmt.Errorf("check %s: failed to parse interval: %w", config.GetName(), err)
		}
		nextRun = time.Now().In(s.location).Add(interval)
	} else {
		scheduleType = ScheduleTypeCron
		nextRun, err = s.calculateNextRun(config.GetSchedule())
		if err != nil {
			return fmt.Errorf("check %s: failed to calculate next run: %w", config.GetName(), err)
		}
	}

	s.checks[config.GetName()] = &ScheduledCheck{
		Config:       config,
		NextRun:      nextRun,
		ScheduleType: scheduleType,
		Interval:     interval,
	}

	logger.Info("Added check %s (enabled: %v, next run: %v)", config.GetName(), config.IsEnabled(), nextRun.Format(time.RFC3339))

	return nil
}

func (s *Scheduler) RemoveCheck(name string) {
	delete(s.checks, name)
}

func (s *Scheduler) Start(interval time.Duration) {
	logger.Info("Scheduler started with ticker interval %v", interval)
	s.ticker = time.NewTicker(interval)
	go s.run()
}

func (s *Scheduler) Stop() {
	logger.Info("Scheduler stopping")
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
	logger.Info("Scheduler stopped")
}

func (s *Scheduler) run() {
	for {
		select {
		case <-s.ticker.C:
			s.tick()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) tick() {
	now := time.Now().In(s.location)
	logger.Debug("Scheduler tick at %v", now.Format(time.RFC3339))

	s.processQueue()

	for name, check := range s.checks {
		if !check.Config.IsEnabled() || check.Running {
			continue
		}

		if now.After(check.NextRun) || now.Equal(check.NextRun) {
			s.executeCheck(name, check)
		}
	}
}

func (s *Scheduler) processQueue() {
	if s.maxConcurrentChecks <= 0 {
		return
	}

	for len(s.queue) > 0 && s.runningChecks < s.maxConcurrentChecks {
		checkConfig := s.queue[0]
		s.queue = s.queue[1:]

		if check, exists := s.checks[checkConfig.GetName()]; exists {
			logger.Info("Processing queued check: %s (running: %d, queued: %d)", checkConfig.GetName(), s.runningChecks, len(s.queue))
			s.executeCheck(checkConfig.GetName(), check)
		}
	}
}

func (s *Scheduler) executeCheck(name string, check *ScheduledCheck) {
	if s.maxConcurrentChecks > 0 && s.runningChecks >= s.maxConcurrentChecks {
		logger.Debug("Queuing check %s (concurrency limit reached: %d)", name, s.maxConcurrentChecks)
		s.queue = append(s.queue, check.Config)
		return
	}

	check.Running = true
	s.runningChecks++
	now := time.Now().In(s.location)
	check.LastRun = &now

	logger.Info("Executing check: %s (next run: %v)", name, check.NextRun.Format(time.RFC3339))

	go func() {
		defer func() {
			check.Running = false
			s.runningChecks--
			now := time.Now().In(s.location)
			if check.ScheduleType == ScheduleTypeInterval && check.Interval > 0 {
				check.NextRun = now.Add(check.Interval)
			} else {
				nextRun, err := s.calculateNextRun(check.Config.GetSchedule())
				if err == nil {
					check.NextRun = nextRun
				}
			}
			check.LastRun = &now
			logger.Debug("Check %s completed (next run: %v)", name, check.NextRun.Format(time.RFC3339))
		}()

		if err := s.executor.Execute(check.Config); err != nil {
			logger.Error("Error executing check %s: %v", name, err)
		}
	}()
}

func (s *Scheduler) calculateNextRun(cronExpr string) (time.Time, error) {
	parsed, err := ParseCronExpression(cronExpr)
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now().In(s.location)
	return parsed.Next(now), nil
}

func parseInterval(interval string) (time.Duration, error) {
	interval = strings.TrimSpace(interval)
	if interval == "" {
		return 0, fmt.Errorf("interval cannot be empty")
	}

	unit := interval[len(interval)-1:]
	valueStr := interval[:len(interval)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid interval value: %s", valueStr)
	}

	if value <= 0 {
		return 0, fmt.Errorf("interval value must be positive: %d", value)
	}

	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid interval unit: %s (must be s, m, h, or d)", unit)
	}
}

func (s *Scheduler) GetCheckStatus(name string) (*ScheduledCheck, bool) {
	check, exists := s.checks[name]
	return check, exists
}

func (s *Scheduler) GetAllChecks() map[string]*ScheduledCheck {
	return s.checks
}
