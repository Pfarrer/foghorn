package daemon

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/pfarrer/foghorn/config"
	"github.com/pfarrer/foghorn/containerimage"
	"github.com/pfarrer/foghorn/scheduler"
)

func TestVerifyImageAvailability_NoChecks(t *testing.T) {
	cfg := &config.Config{
		Checks: []config.CheckConfig{},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error with no checks, got: %v", err)
	}
}

func TestVerifyImageAvailability_NoEnabledChecks(t *testing.T) {
	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "disabled-check",
				Image:   "test/image:1.0.0",
				Enabled: false,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error with no enabled checks, got: %v", err)
	}
}

func TestVerifyImageAvailability_MissingImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "test-check",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	expected := "Error: The following Docker images are not available locally"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error message to contain '%s', got: %v", expected, err)
	}

	expectedPull := "docker pull this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expectedPull) {
		t.Errorf("Expected error message to contain '%s', got: %v", expectedPull, err)
	}
}

func TestVerifyImageAvailability_ExistingImage(t *testing.T) {
	image := findLocalSemverImage(t)

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "alpine-check",
				Image:   image,
				Enabled: true,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err != nil {
		t.Errorf("Expected no error for existing image, got: %v", err)
	}
}

func TestVerifyImageAvailability_MultipleChecksSameImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "check1",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	if err != nil && !containsString(err.Error(), "(required by: check1, check2)") &&
		!containsString(err.Error(), "(required by: check2, check1)") {
		t.Errorf("Expected error to list both check names, got: %v", err)
	}
}

func TestVerifyImageAvailability_MultipleImages(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "check1",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
			{
				Name:    "check2",
				Image:   "another-missing-image-67890:2.3.4",
				Enabled: true,
			},
		},
	}

	err = verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing images, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	expected2 := "another-missing-image-67890:2.3.4"
	if err != nil && !containsString(err.Error(), expected2) {
		t.Errorf("Expected error to contain '%s', got: %v", expected2, err)
	}
}

func TestVerifyImageAvailability_MixedMissingAndExisting(t *testing.T) {
	image := findLocalSemverImage(t)

	cfg := &config.Config{
		Checks: []config.CheckConfig{
			{
				Name:    "existing-check",
				Image:   image,
				Enabled: true,
			},
			{
				Name:    "missing-check",
				Image:   "this-image-definitely-does-not-exist-12345:1.2.3",
				Enabled: true,
			},
		},
	}

	err := verifyImageAvailabilityFn(cfg)
	if err == nil {
		t.Error("Expected error for missing image, got nil")
	}

	expected := "this-image-definitely-does-not-exist-12345:1.2.3"
	if err != nil && !containsString(err.Error(), expected) {
		t.Errorf("Expected error to contain '%s', got: %v", expected, err)
	}

	if err != nil && containsString(err.Error(), image) {
		t.Errorf("Error should not mention existing image %s, got: %v", image, err)
	}
}

func findLocalSemverImage(t *testing.T) string {
	t.Helper()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skipf("Skipping test: cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), image.ListOptions{})
	if err != nil {
		t.Skipf("Skipping test: cannot list docker images: %v", err)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			ref, err := containerimage.ParseReference(tag)
			if err != nil {
				continue
			}
			if ref.Selector.Kind != containerimage.SelectorFull {
				continue
			}
			return tag
		}
	}

	t.Skipf("Skipping test: no local images with semantic version tags found")
	return ""
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type oneShotMockExecutor struct {
	mu        sync.Mutex
	executed  []string
	statusMap map[string]string
	callback  func(checkName string, status string, duration time.Duration, startTime time.Time, err error)
	execCount atomic.Int32
}

func (e *oneShotMockExecutor) Execute(check scheduler.CheckConfig) error {
	e.mu.Lock()
	e.executed = append(e.executed, check.GetName())
	status := "pass"
	if s, ok := e.statusMap[check.GetName()]; ok {
		status = s
	}
	e.mu.Unlock()

	e.execCount.Add(1)

	if e.callback != nil {
		e.callback(check.GetName(), status, 10*time.Millisecond, time.Time{}, nil)
	}
	return nil
}

func (e *oneShotMockExecutor) SetResultCallback(callback func(checkName string, status string, duration time.Duration, startTime time.Time, err error)) {
	e.callback = callback
}

func (e *oneShotMockExecutor) getExecuted() []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make([]string, len(e.executed))
	copy(result, e.executed)
	return result
}

func buildOneShotScheduler(exec scheduler.CheckExecutor, checks []config.CheckConfig) *scheduler.Scheduler {
	sched := scheduler.NewScheduler(exec, time.UTC, 0)
	for i := range checks {
		adapter := scheduler.NewConfigAdapter(&checks[i])
		sched.AddCheck(adapter)
	}
	return sched
}

func TestRunOneShot_FlagEnablesMode(t *testing.T) {
	exec := &oneShotMockExecutor{
		statusMap: map[string]string{
			"check-a": "pass",
		},
	}

	checks := []config.CheckConfig{
		{Name: "check-a", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
	}
	sched := buildOneShotScheduler(exec, checks)

	exitCode := runOneShot(sched, exec, 0)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	executed := exec.getExecuted()
	if len(executed) != 1 || executed[0] != "check-a" {
		t.Errorf("expected [check-a] executed, got %v", executed)
	}
}

func TestRunOneShot_EachCheckRunsOnce(t *testing.T) {
	exec := &oneShotMockExecutor{
		statusMap: map[string]string{
			"check-1": "pass",
			"check-2": "pass",
			"check-3": "pass",
			"check-4": "pass",
			"check-5": "pass",
		},
	}

	checks := []config.CheckConfig{
		{Name: "check-1", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "check-2", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "check-3", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "check-4", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "check-5", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
	}
	sched := buildOneShotScheduler(exec, checks)

	exitCode := runOneShot(sched, exec, 0)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if count := exec.execCount.Load(); count != 5 {
		t.Errorf("expected 5 executions, got %d", count)
	}

	executed := exec.getExecuted()
	if len(executed) != 5 {
		t.Errorf("expected 5 unique executions, got %v", executed)
	}
}

func TestRunOneShot_ExitCode0WhenAllPassOrWarn(t *testing.T) {
	exec := &oneShotMockExecutor{
		statusMap: map[string]string{
			"pass-check": "pass",
			"warn-check": "warn",
		},
	}

	checks := []config.CheckConfig{
		{Name: "pass-check", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "warn-check", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
	}
	sched := buildOneShotScheduler(exec, checks)

	exitCode := runOneShot(sched, exec, 0)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 when all pass/warn, got %d", exitCode)
	}
}

func TestRunOneShot_ExitCode1WhenAnyFailOrUnknown(t *testing.T) {
	tests := []struct {
		name      string
		statusMap map[string]string
		wantCode  int
	}{
		{
			name:      "fail status",
			statusMap: map[string]string{"check-a": "fail"},
			wantCode:  1,
		},
		{
			name:      "unknown status",
			statusMap: map[string]string{"check-a": "unknown"},
			wantCode:  1,
		},
		{
			name:      "error status",
			statusMap: map[string]string{"check-a": "error"},
			wantCode:  1,
		},
		{
			name:      "mixed pass and fail",
			statusMap: map[string]string{"pass-check": "pass", "fail-check": "fail"},
			wantCode:  1,
		},
		{
			name:      "mixed pass and unknown",
			statusMap: map[string]string{"pass-check": "pass", "unknown-check": "unknown"},
			wantCode:  1,
		},
		{
			name:      "mixed pass and error",
			statusMap: map[string]string{"pass-check": "pass", "error-check": "error"},
			wantCode:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := &oneShotMockExecutor{statusMap: tt.statusMap}

			var checks []config.CheckConfig
			for name := range tt.statusMap {
				checks = append(checks, config.CheckConfig{
					Name: name, Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"},
				})
			}
			sched := buildOneShotScheduler(exec, checks)

			exitCode := runOneShot(sched, exec, 0)

			if exitCode != tt.wantCode {
				t.Errorf("expected exit code %d, got %d", tt.wantCode, exitCode)
			}
		})
	}
}

func TestRunOneShot_DisabledChecksSkipped(t *testing.T) {
	exec := &oneShotMockExecutor{
		statusMap: map[string]string{
			"enabled-check": "pass",
		},
	}

	checks := []config.CheckConfig{
		{Name: "enabled-check", Enabled: true, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
		{Name: "disabled-check", Enabled: false, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
	}
	sched := buildOneShotScheduler(exec, checks)

	exitCode := runOneShot(sched, exec, 0)

	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	if count := exec.execCount.Load(); count != 1 {
		t.Errorf("expected 1 execution (disabled check skipped), got %d", count)
	}
}

func TestRunOneShot_NoEnabledChecks(t *testing.T) {
	exec := &oneShotMockExecutor{statusMap: map[string]string{}}

	checks := []config.CheckConfig{
		{Name: "disabled-check", Enabled: false, Schedule: config.Schedule{Cron: "*/5 * * * *"}},
	}
	sched := buildOneShotScheduler(exec, checks)

	exitCode := runOneShot(sched, exec, 0)

	if exitCode != 0 {
		t.Errorf("expected exit code 0 with no enabled checks, got %d", exitCode)
	}

	if count := exec.execCount.Load(); count != 0 {
		t.Errorf("expected 0 executions, got %d", count)
	}
}
