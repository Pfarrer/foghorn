package scheduler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	SetResultCallback(callback func(checkName string, status string, duration time.Duration))
}

type ScheduledCheck struct {
	Config       CheckConfig
	NextRun      time.Time
	LastRun      *time.Time
	LastStatus   string
	LastDuration time.Duration
	Running      bool
	ScheduleType ScheduleType
	Interval     time.Duration
	IsQueued     bool
	History      []CheckHistoryEntry
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
	startTime           time.Time
	mu                  sync.RWMutex
	resultLogger        ResultLogger
}

type ResultLogger interface {
	RecordResult(checkName string, status string, duration time.Duration, completedAt time.Time) error
}

type CheckHistoryEntry struct {
	Status      string
	CompletedAt time.Time
}

const maxHistoryEntries = 10

func NewScheduler(executor CheckExecutor, location *time.Location, maxConcurrentChecks int) *Scheduler {
	if location == nil {
		location = time.UTC
	}
	s := &Scheduler{
		checks:              make(map[string]*ScheduledCheck),
		executor:            executor,
		stopChan:            make(chan struct{}),
		location:            location,
		maxConcurrentChecks: maxConcurrentChecks,
		queue:               make([]CheckConfig, 0),
		startTime:           time.Now(),
	}

	executor.SetResultCallback(s.handleCheckResult)

	return s
}

func (s *Scheduler) SetResultLogger(logger ResultLogger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resultLogger = logger
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
		nextRun = time.Now().In(s.location)
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
		LastStatus:   "unknown",
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

	s.processQueue()

	var due []struct {
		name  string
		check *ScheduledCheck
	}

	s.mu.RLock()
	for name, check := range s.checks {
		if !check.Config.IsEnabled() || check.Running {
			continue
		}
		if now.After(check.NextRun) || now.Equal(check.NextRun) {
			due = append(due, struct {
				name  string
				check *ScheduledCheck
			}{name: name, check: check})
		}
	}
	s.mu.RUnlock()

	if len(due) > 1 {
		sort.Slice(due, func(i, j int) bool {
			pi := s.priorityDuration(due[i].check, now)
			pj := s.priorityDuration(due[j].check, now)
			if pi == pj {
				return due[i].name < due[j].name
			}
			return pi > pj
		})
	}

	for _, item := range due {
		s.executeCheck(item.name, item.check)
	}
}

func (s *Scheduler) processQueue() {
	if s.maxConcurrentChecks <= 0 {
		return
	}

	for {
		var checkConfig CheckConfig
		s.mu.Lock()
		if len(s.queue) == 0 || s.runningChecks >= s.maxConcurrentChecks {
			s.mu.Unlock()
			break
		}
		s.sortQueueLocked(time.Now().In(s.location))
		checkConfig = s.queue[0]
		s.queue = s.queue[1:]
		s.mu.Unlock()

		if check, exists := s.GetCheckStatus(checkConfig.GetName()); exists {
			logger.Info("Processing queued check: %s (running: %d, queued: %d)", checkConfig.GetName(), s.runningChecks, len(s.queue))
			s.mu.Lock()
			check.IsQueued = false
			s.mu.Unlock()
			s.executeCheck(checkConfig.GetName(), check)
		}
	}

	s.mu.Lock()
	for name, check := range s.checks {
		isInQueue := false
		for _, queuedCheck := range s.queue {
			if queuedCheck.GetName() == name {
				isInQueue = true
				break
			}
		}
		check.IsQueued = isInQueue
	}
	s.mu.Unlock()
}

func (s *Scheduler) executeCheck(name string, check *ScheduledCheck) {
	s.mu.Lock()
	if s.maxConcurrentChecks > 0 && s.runningChecks >= s.maxConcurrentChecks {
		logger.Debug("Queuing check %s (concurrency limit reached: %d)", name, s.maxConcurrentChecks)
		s.queue = append(s.queue, check.Config)
		s.sortQueueLocked(time.Now().In(s.location))
		check.IsQueued = true
		s.mu.Unlock()
		return
	}

	check.Running = true
	check.IsQueued = false
	s.runningChecks++
	now := time.Now().In(s.location)
	check.LastRun = &now
	s.mu.Unlock()

	logger.Info("Executing check: %s (next run: %v)", name, check.NextRun.Format(time.RFC3339))

	startTime := time.Now()
	go func() {
		defer func() {
			s.mu.Lock()
			check.Running = false
			s.runningChecks--
			now := time.Now().In(s.location)
			check.LastDuration = now.Sub(startTime)
			if check.ScheduleType == ScheduleTypeInterval && check.Interval > 0 {
				check.NextRun = now.Add(check.Interval)
			} else {
				nextRun, err := s.calculateNextRun(check.Config.GetSchedule())
				if err == nil {
					check.NextRun = nextRun
				}
			}
			check.LastRun = &now
			s.mu.Unlock()
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

func (s *Scheduler) priorityDuration(check *ScheduledCheck, now time.Time) time.Duration {
	if check == nil {
		return 0
	}
	if check.ScheduleType == ScheduleTypeInterval && check.Interval > 0 {
		return check.Interval
	}
	if check.ScheduleType == ScheduleTypeCron {
		if next, err := s.calculateNextRun(check.Config.GetSchedule()); err == nil {
			return next.Sub(now)
		}
	}
	if !check.NextRun.IsZero() {
		return check.NextRun.Sub(now)
	}
	return 0
}

func (s *Scheduler) sortQueueLocked(now time.Time) {
	if len(s.queue) < 2 {
		return
	}
	sort.Slice(s.queue, func(i, j int) bool {
		ci := s.checks[s.queue[i].GetName()]
		cj := s.checks[s.queue[j].GetName()]
		pi := s.priorityDuration(ci, now)
		pj := s.priorityDuration(cj, now)
		if pi == pj {
			return s.queue[i].GetName() < s.queue[j].GetName()
		}
		return pi > pj
	})
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	check, exists := s.checks[name]
	return check, exists
}

func (s *Scheduler) GetAllChecks() map[string]*ScheduledCheck {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]*ScheduledCheck, len(s.checks))
	for k, v := range s.checks {
		result[k] = v
	}
	return result
}

func (s *Scheduler) GetStartTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.startTime
}

func (s *Scheduler) GetCounts() (total, running, queued, pass, fail, warn int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total = len(s.checks)
	running = s.runningChecks
	queued = len(s.queue)

	for _, check := range s.checks {
		switch check.LastStatus {
		case "pass":
			pass++
		case "fail":
			fail++
		case "warn":
			warn++
		}
	}

	return
}

func (s *Scheduler) handleCheckResult(checkName string, status string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	completedAt := time.Now().In(s.location)
	if check, exists := s.checks[checkName]; exists {
		check.LastStatus = status
		check.LastDuration = duration
		check.History = trimHistory(append(check.History, CheckHistoryEntry{
			Status:      status,
			CompletedAt: completedAt,
		}))
	}

	if s.resultLogger != nil {
		if err := s.resultLogger.RecordResult(checkName, status, duration, completedAt); err != nil {
			logger.Error("Failed to persist state for %s: %v", checkName, err)
		}
	}
}

type CheckState struct {
	LastStatus   string
	LastDuration time.Duration
	LastRun      time.Time
	History      []CheckHistoryEntry
}

func (s *Scheduler) ApplyState(states map[string]CheckState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, state := range states {
		check, exists := s.checks[name]
		if !exists {
			continue
		}
		if state.LastStatus != "" {
			check.LastStatus = state.LastStatus
		}
		if state.LastDuration > 0 {
			check.LastDuration = state.LastDuration
		}
		if !state.LastRun.IsZero() {
			lastRun := state.LastRun
			check.LastRun = &lastRun
			if check.ScheduleType == ScheduleTypeInterval && check.Interval > 0 {
				check.NextRun = lastRun.Add(check.Interval)
			}
		}
		if len(state.History) > 0 {
			check.History = trimHistory(state.History)
		}
	}
}

func trimHistory(entries []CheckHistoryEntry) []CheckHistoryEntry {
	if len(entries) == 0 {
		return nil
	}
	if len(entries) > maxHistoryEntries {
		entries = entries[len(entries)-maxHistoryEntries:]
	}
	out := make([]CheckHistoryEntry, len(entries))
	copy(out, entries)
	return out
}
