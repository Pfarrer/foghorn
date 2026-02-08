package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/scheduler"
)

type stubConfig struct {
	name     string
	schedule string
	enabled  bool
}

func (s *stubConfig) GetName() string     { return s.name }
func (s *stubConfig) GetSchedule() string { return s.schedule }
func (s *stubConfig) IsEnabled() bool     { return s.enabled }

type stubExecutor struct{}

func (s *stubExecutor) Execute(check scheduler.CheckConfig) error { return nil }
func (s *stubExecutor) SetResultCallback(callback func(string, string, time.Duration)) {
}

func TestCheckHeaderColumns(t *testing.T) {
	sched := scheduler.NewScheduler(&stubExecutor{}, time.UTC, 0)
	model := NewModel(sched, "info")
	model.width = 120

	styles := newStyles(model.width)
	header := model.renderCheckHeader(12, styles)

	if strings.Contains(header, "  Status  ") {
		t.Fatalf("header should not include status column: %q", header)
	}
	if !strings.Contains(header, "Last Status") {
		t.Fatalf("header missing Last Status column: %q", header)
	}
	if !strings.Contains(header, "Status History") {
		t.Fatalf("header missing Status History column: %q", header)
	}
}

func TestHistorySymbolsRender(t *testing.T) {
	styles := newStyles(120)
	entries := []scheduler.CheckHistoryEntry{
		{Status: "pass", CompletedAt: time.Now().Add(-2 * time.Minute)},
		{Status: "fail", CompletedAt: time.Now().Add(-1 * time.Minute)},
	}

	history := formatHistorySymbols(entries, 10, styles)
	if history == "-" {
		t.Fatalf("history should render symbols, got %q", history)
	}
}
