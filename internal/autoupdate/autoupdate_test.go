package autoupdate

import (
	"context"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/config"
)

func TestAutoUpdater_SkipsDisabledChecks(t *testing.T) {
	checks := []config.CheckConfig{
		{Name: "disabled-check", Image: "test/image:1.0.0", Enabled: false},
		{Name: "enabled-check", Image: "test/image:2.0.0", Enabled: true},
	}

	u := &AutoUpdater{checks: checks}

	if len(u.checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(u.checks))
	}
}

func TestAutoUpdater_NewAutoUpdater(t *testing.T) {
	checks := []config.CheckConfig{
		{Name: "test-check", Image: "test/image:1.0.0", Enabled: true},
	}

	u := NewAutoUpdater(nil, checks)

	if u == nil {
		t.Fatal("expected non-nil AutoUpdater")
	}
	if len(u.checks) != 1 {
		t.Errorf("expected 1 check, got %d", len(u.checks))
	}
}

func TestAutoUpdater_Run_ContextCancellation(t *testing.T) {
	checks := []config.CheckConfig{}
	u := NewAutoUpdater(nil, checks)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		u.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not exit on cancelled context")
	}
}

func TestCheckNeedsResolution(t *testing.T) {
	tests := []struct {
		image    string
		expected bool
	}{
		{"test/image:1.0.0", false},
		{"test/image:1", true},
		{"test/image:1.2", true},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			result := checkNeedsResolution(tt.image)
			if result != tt.expected {
				t.Errorf("checkNeedsResolution(%q) = %v, want %v", tt.image, result, tt.expected)
			}
		})
	}
}
