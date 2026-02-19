package tui

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pfarrer/foghorn/internal/statusapi"
	"github.com/pfarrer/foghorn/scheduler"
)

type remoteStatusReader struct {
	client   *statusapi.Client
	mu       sync.RWMutex
	started  time.Time
	counts   scheduler.SnapshotCounts
	checks   map[string]*scheduler.ScheduledCheck
	snapshot time.Time
}

func newRemoteStatusReader(statusURL string) (*remoteStatusReader, error) {
	url := strings.TrimSpace(statusURL)
	if url == "" {
		url = statusapi.DefaultBaseURL
	}
	return &remoteStatusReader{
		client: statusapi.NewClient(url),
		checks: make(map[string]*scheduler.ScheduledCheck),
	}, nil
}

func (r *remoteStatusReader) Refresh(ctx context.Context) error {
	s, err := r.client.GetStatus(ctx)
	if err != nil {
		return err
	}

	checks := make(map[string]*scheduler.ScheduledCheck, len(s.Checks))
	for name, check := range s.Checks {
		history := make([]scheduler.CheckHistoryEntry, len(check.History))
		copy(history, check.History)
		duration := time.Duration(check.LastDurationMs) * time.Millisecond
		checks[name] = &scheduler.ScheduledCheck{
			NextRun:      check.NextRun,
			LastRun:      copyTime(check.LastRun),
			LastStatus:   check.LastStatus,
			LastDuration: duration,
			Running:      check.Running,
			ScheduleType: check.ScheduleType,
			IsQueued:     check.Queued,
			History:      history,
		}
	}

	r.mu.Lock()
	r.started = s.StartedAt
	r.counts = s.Counts
	r.checks = checks
	r.snapshot = s.GeneratedAt
	r.mu.Unlock()
	return nil
}

func (r *remoteStatusReader) GetStartTime() time.Time {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.started
}

func (r *remoteStatusReader) GetCounts() (total, running, queued, pass, fail, warn int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counts.Total, r.counts.Running, r.counts.Queued, r.counts.Pass, r.counts.Fail, r.counts.Warn
}

func (r *remoteStatusReader) GetAllChecks() map[string]*scheduler.ScheduledCheck {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make(map[string]*scheduler.ScheduledCheck, len(r.checks))
	for name, check := range r.checks {
		history := make([]scheduler.CheckHistoryEntry, len(check.History))
		copy(history, check.History)
		out[name] = &scheduler.ScheduledCheck{
			NextRun:      check.NextRun,
			LastRun:      copyTime(check.LastRun),
			LastStatus:   check.LastStatus,
			LastDuration: check.LastDuration,
			Running:      check.Running,
			ScheduleType: check.ScheduleType,
			IsQueued:     check.IsQueued,
			History:      history,
		}
	}
	return out
}

func copyTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}
