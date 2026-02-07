package scheduler

import (
	"fmt"
	"testing"
	"time"
)

type MockCheckConfig struct {
	name     string
	schedule string
	enabled  bool
}

func (m *MockCheckConfig) GetName() string {
	return m.name
}

func (m *MockCheckConfig) GetSchedule() string {
	return m.schedule
}

func (m *MockCheckConfig) IsEnabled() bool {
	return m.enabled
}

type IntervalMockCheckConfig struct {
	name     string
	schedule string
	enabled  bool
	interval string
}

func (m *IntervalMockCheckConfig) GetName() string {
	return m.name
}

func (m *IntervalMockCheckConfig) GetSchedule() string {
	return m.interval
}

func (m *IntervalMockCheckConfig) IsEnabled() bool {
	return m.enabled
}

func (m *IntervalMockCheckConfig) GetScheduleType() ScheduleType {
	return ScheduleTypeInterval
}

func (m *IntervalMockCheckConfig) GetInterval() string {
	return m.interval
}

type MockExecutor struct {
	executed []string
	callback func(checkName string, status string, duration time.Duration)
}

func (m *MockExecutor) Execute(check CheckConfig) error {
	m.executed = append(m.executed, check.GetName())
	return nil
}

func (m *MockExecutor) SetResultCallback(callback func(checkName string, status string, duration time.Duration)) {
	m.callback = callback
}

type SlowExecutor struct {
	executed []string
	blocker  chan struct{}
	callback func(checkName string, status string, duration time.Duration)
}

func (m *SlowExecutor) Execute(check CheckConfig) error {
	m.executed = append(m.executed, check.GetName())
	<-m.blocker
	return nil
}

func (m *SlowExecutor) SetResultCallback(callback func(checkName string, status string, duration time.Duration)) {
	m.callback = callback
}

func TestParseCronExpression(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr bool
	}{
		{
			name:    "simple expression",
			expr:    "0 12 * * *",
			wantErr: false,
		},
		{
			name:    "wildcard all",
			expr:    "* * * * *",
			wantErr: false,
		},
		{
			name:    "list values",
			expr:    "0,15,30,45 * * * *",
			wantErr: false,
		},
		{
			name:    "range",
			expr:    "0 9-17 * * *",
			wantErr: false,
		},
		{
			name:    "step",
			expr:    "*/5 * * * *",
			wantErr: false,
		},
		{
			name:    "combined",
			expr:    "0,30 9-17 * * 1-5",
			wantErr: false,
		},
		{
			name:    "invalid - too many fields",
			expr:    "* * * * * *",
			wantErr: true,
		},
		{
			name:    "invalid - too few fields",
			expr:    "* * * *",
			wantErr: true,
		},
		{
			name:    "invalid minute",
			expr:    "60 * * * *",
			wantErr: true,
		},
		{
			name:    "invalid hour",
			expr:    "* 24 * * *",
			wantErr: true,
		},
		{
			name:    "invalid month",
			expr:    "* * * 13 *",
			wantErr: true,
		},
		{
			name:    "invalid day of week",
			expr:    "* * * * 7",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCronExpression(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCronExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCronExpressionNext(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		baseTime time.Time
		wantNext time.Time
	}{
		{
			name:     "every minute",
			expr:     "* * * * *",
			baseTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantNext: time.Date(2024, 1, 1, 12, 1, 0, 0, time.UTC),
		},
		{
			name:     "every 5 minutes",
			expr:     "*/5 * * * *",
			baseTime: time.Date(2024, 1, 1, 12, 2, 0, 0, time.UTC),
			wantNext: time.Date(2024, 1, 1, 12, 5, 0, 0, time.UTC),
		},
		{
			name:     "specific minute",
			expr:     "30 * * * *",
			baseTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantNext: time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			name:     "daily at midnight",
			expr:     "0 0 * * *",
			baseTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantNext: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron, err := ParseCronExpression(tt.expr)
			if err != nil {
				t.Fatalf("Failed to parse cron expression: %v", err)
			}

			next := cron.Next(tt.baseTime)
			if !next.Equal(tt.wantNext) {
				t.Errorf("Next() = %v, want %v", next, tt.wantNext)
			}
		})
	}
}

func TestSchedulerAddCheck(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "test-check",
		schedule: "*/5 * * * *",
		enabled:  true,
	}

	err := scheduler.AddCheck(check)
	if err != nil {
		t.Fatalf("AddCheck() error = %v", err)
	}

	checkStatus, exists := scheduler.GetCheckStatus("test-check")
	if !exists {
		t.Error("Check should exist after being added")
	}

	if checkStatus.Config.GetName() != "test-check" {
		t.Errorf("Check name should be 'test-check', got %s", checkStatus.Config.GetName())
	}
}

func TestSchedulerAddCheckInvalidSchedule(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "test-check",
		schedule: "invalid",
		enabled:  true,
	}

	err := scheduler.AddCheck(check)
	if err == nil {
		t.Error("AddCheck() should return error for invalid schedule")
	}
}

func TestSchedulerRemoveCheck(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "test-check",
		schedule: "*/5 * * * *",
		enabled:  true,
	}

	scheduler.AddCheck(check)
	scheduler.RemoveCheck("test-check")

	_, exists := scheduler.GetCheckStatus("test-check")
	if exists {
		t.Error("Check should not exist after being removed")
	}
}

func TestSchedulerExecution(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "test-check",
		schedule: "* * * * *",
		enabled:  true,
	}

	err := scheduler.AddCheck(check)
	if err != nil {
		t.Fatalf("AddCheck() error = %v", err)
	}

	checkStatus, exists := scheduler.GetCheckStatus("test-check")
	if !exists {
		t.Fatal("Check should exist")
	}

	if checkStatus.NextRun.IsZero() {
		t.Error("Next run time should be set")
	}

	if checkStatus.Running {
		t.Error("Check should not be running initially")
	}
}

func TestSchedulerDisabledCheck(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "disabled-check",
		schedule: "* * * * *",
		enabled:  false,
	}

	err := scheduler.AddCheck(check)
	if err != nil {
		t.Fatalf("AddCheck() error = %v", err)
	}

	scheduler.Start(10 * time.Millisecond)
	time.Sleep(150 * time.Millisecond)
	scheduler.Stop()

	if len(executor.executed) > 0 {
		t.Error("Disabled check should not be executed")
	}
}

func TestCronFieldMatches(t *testing.T) {
	tests := []struct {
		name  string
		field CronField
		value int
		want  bool
	}{
		{
			name: "matches within range",
			field: CronField{
				min:    0,
				max:    59,
				values: map[int]bool{0: true, 5: true, 10: true},
			},
			value: 5,
			want:  true,
		},
		{
			name: "does not match",
			field: CronField{
				min:    0,
				max:    59,
				values: map[int]bool{0: true, 5: true, 10: true},
			},
			value: 15,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.matches(tt.value); got != tt.want {
				t.Errorf("CronField.matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeZones(t *testing.T) {
	tests := []struct {
		name     string
		location *time.Location
		expr     string
		wantErr  bool
	}{
		{
			name:     "UTC timezone",
			location: time.UTC,
			expr:     "0 12 * * *",
			wantErr:  false,
		},
		{
			name:     "America/New_York timezone",
			location: time.FixedZone("EST", -5*3600),
			expr:     "0 12 * * *",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &MockExecutor{}
			scheduler := NewScheduler(executor, tt.location, 0)

			check := &MockCheckConfig{
				name:     "test-check",
				schedule: tt.expr,
				enabled:  true,
			}

			err := scheduler.AddCheck(check)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConcurrencyLimit(t *testing.T) {
	e := &SlowExecutor{
		executed: []string{},
		blocker:  make(chan struct{}),
	}

	scheduler := NewScheduler(e, time.UTC, 2)

	now := time.Now()

	for i := 1; i <= 5; i++ {
		check := &MockCheckConfig{
			name:     fmt.Sprintf("check-%d", i),
			schedule: "* * * * *",
			enabled:  true,
		}
		err := scheduler.AddCheck(check)
		if err != nil {
			t.Fatalf("AddCheck() error = %v", err)
		}
		checkStatus, _ := scheduler.GetCheckStatus(check.name)
		checkStatus.NextRun = now
	}

	scheduler.Start(10 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)

	allChecks := scheduler.GetAllChecks()
	runningCount := 0
	for _, check := range allChecks {
		if check.Running {
			runningCount++
		}
	}

	if runningCount > 2 {
		t.Errorf("Expected max 2 running checks, got %d", runningCount)
	}

	scheduler.Stop()
	close(e.blocker)
}

func TestUnlimitedConcurrency(t *testing.T) {
	e := &SlowExecutor{
		executed: []string{},
		blocker:  make(chan struct{}),
	}

	scheduler := NewScheduler(e, time.UTC, 0)

	now := time.Now()

	for i := 1; i <= 5; i++ {
		check := &MockCheckConfig{
			name:     fmt.Sprintf("check-%d", i),
			schedule: "* * * * *",
			enabled:  true,
		}
		err := scheduler.AddCheck(check)
		if err != nil {
			t.Fatalf("AddCheck() error = %v", err)
		}
		checkStatus, _ := scheduler.GetCheckStatus(check.name)
		checkStatus.NextRun = now
	}

	scheduler.Start(10 * time.Millisecond)
	time.Sleep(100 * time.Millisecond)

	allChecks := scheduler.GetAllChecks()
	runningCount := 0
	for _, check := range allChecks {
		if check.Running {
			runningCount++
		}
	}

	if runningCount != 5 {
		t.Errorf("Expected 5 running checks with unlimited concurrency, got %d", runningCount)
	}

	for i := 0; i < 5; i++ {
		e.blocker <- struct{}{}
	}

	scheduler.Stop()
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		wantErr  bool
		wantDur  time.Duration
	}{
		{
			name:     "seconds",
			interval: "30s",
			wantErr:  false,
			wantDur:  30 * time.Second,
		},
		{
			name:     "minutes",
			interval: "5m",
			wantErr:  false,
			wantDur:  5 * time.Minute,
		},
		{
			name:     "hours",
			interval: "2h",
			wantErr:  false,
			wantDur:  2 * time.Hour,
		},
		{
			name:     "days",
			interval: "1d",
			wantErr:  false,
			wantDur:  24 * time.Hour,
		},
		{
			name:     "empty interval",
			interval: "",
			wantErr:  true,
		},
		{
			name:     "invalid unit",
			interval: "5x",
			wantErr:  true,
		},
		{
			name:     "missing unit",
			interval: "5",
			wantErr:  true,
		},
		{
			name:     "negative value",
			interval: "-5m",
			wantErr:  true,
		},
		{
			name:     "zero value",
			interval: "0s",
			wantErr:  true,
		},
		{
			name:     "non-numeric value",
			interval: "abc",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseInterval(tt.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseInterval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantDur {
				t.Errorf("parseInterval() = %v, want %v", got, tt.wantDur)
			}
		})
	}
}

func TestIntervalBasedScheduling(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &IntervalMockCheckConfig{
		name:     "interval-check",
		schedule: "30s",
		enabled:  true,
		interval: "30s",
	}

	err := scheduler.AddCheck(check)
	if err != nil {
		t.Fatalf("AddCheck() error = %v", err)
	}

	checkStatus, exists := scheduler.GetCheckStatus("interval-check")
	if !exists {
		t.Fatal("Check should exist")
	}

	if checkStatus.ScheduleType != ScheduleTypeInterval {
		t.Errorf("Expected schedule type interval, got %v", checkStatus.ScheduleType)
	}

	if checkStatus.Interval != 30*time.Second {
		t.Errorf("Expected interval 30s, got %v", checkStatus.Interval)
	}

	now := time.Now().In(time.UTC)
	if checkStatus.NextRun.Before(now) || checkStatus.NextRun.After(now.Add(31*time.Second)) {
		t.Errorf("Next run time should be approximately 30 seconds from now, got %v", checkStatus.NextRun)
	}
}

func TestMixedCronAndIntervalScheduling(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	cronCheck := &MockCheckConfig{
		name:     "cron-check",
		schedule: "* * * * *",
		enabled:  true,
	}

	intervalCheck := &IntervalMockCheckConfig{
		name:     "interval-check",
		schedule: "1m",
		enabled:  true,
		interval: "1m",
	}

	err := scheduler.AddCheck(cronCheck)
	if err != nil {
		t.Fatalf("AddCheck() for cron check error = %v", err)
	}

	err = scheduler.AddCheck(intervalCheck)
	if err != nil {
		t.Fatalf("AddCheck() for interval check error = %v", err)
	}

	cronStatus, _ := scheduler.GetCheckStatus("cron-check")
	if cronStatus.ScheduleType != ScheduleTypeCron {
		t.Errorf("Expected cron check to have schedule type cron, got %v", cronStatus.ScheduleType)
	}

	intervalStatus, _ := scheduler.GetCheckStatus("interval-check")
	if intervalStatus.ScheduleType != ScheduleTypeInterval {
		t.Errorf("Expected interval check to have schedule type interval, got %v", intervalStatus.ScheduleType)
	}
}

func TestApplyStateUpdatesIntervalNextRun(t *testing.T) {
	executor := &MockExecutor{}
	scheduler := NewScheduler(executor, time.UTC, 0)

	check := &IntervalMockCheckConfig{
		name:     "interval-check",
		schedule: "10s",
		enabled:  true,
		interval: "10s",
	}

	if err := scheduler.AddCheck(check); err != nil {
		t.Fatalf("failed to add check: %v", err)
	}

	lastRun := time.Now().UTC().Add(-5 * time.Minute)
	scheduler.ApplyState(map[string]CheckState{
		"interval-check": {
			LastStatus:   "pass",
			LastDuration: 2 * time.Second,
			LastRun:      lastRun,
		},
	})

	status, ok := scheduler.GetCheckStatus("interval-check")
	if !ok {
		t.Fatalf("check not found after apply state")
	}
	expectedNext := lastRun.Add(10 * time.Second)
	if !status.NextRun.Equal(expectedNext) {
		t.Fatalf("expected next run %v, got %v", expectedNext, status.NextRun)
	}
	if status.LastStatus != "pass" {
		t.Fatalf("expected last status pass, got %s", status.LastStatus)
	}
}
