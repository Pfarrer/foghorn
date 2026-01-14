package scheduler

import (
	"fmt"
	"time"
)

type CheckConfig interface {
	GetName() string
	GetSchedule() string
	IsEnabled() bool
}

type CheckExecutor interface {
	Execute(check CheckConfig) error
}

type ScheduledCheck struct {
	Config  CheckConfig
	NextRun time.Time
	LastRun *time.Time
	Running bool
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

	nextRun, err := s.calculateNextRun(config.GetSchedule())
	if err != nil {
		return fmt.Errorf("check %s: failed to calculate next run: %w", config.GetName(), err)
	}

	s.checks[config.GetName()] = &ScheduledCheck{
		Config:  config,
		NextRun: nextRun,
	}

	return nil
}

func (s *Scheduler) RemoveCheck(name string) {
	delete(s.checks, name)
}

func (s *Scheduler) Start(interval time.Duration) {
	s.ticker = time.NewTicker(interval)
	go s.run()
}

func (s *Scheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	close(s.stopChan)
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
			s.executeCheck(checkConfig.GetName(), check)
		}
	}
}

func (s *Scheduler) executeCheck(name string, check *ScheduledCheck) {
	if s.maxConcurrentChecks > 0 && s.runningChecks >= s.maxConcurrentChecks {
		s.queue = append(s.queue, check.Config)
		return
	}

	check.Running = true
	s.runningChecks++
	now := time.Now().In(s.location)
	check.LastRun = &now

	go func() {
		defer func() {
			check.Running = false
			s.runningChecks--
			nextRun, err := s.calculateNextRun(check.Config.GetSchedule())
			if err == nil {
				check.NextRun = nextRun
			}
		}()

		if err := s.executor.Execute(check.Config); err != nil {
			fmt.Printf("Error executing check %s: %v\n", name, err)
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

func (s *Scheduler) GetCheckStatus(name string) (*ScheduledCheck, bool) {
	check, exists := s.checks[name]
	return check, exists
}

func (s *Scheduler) GetAllChecks() map[string]*ScheduledCheck {
	return s.checks
}
