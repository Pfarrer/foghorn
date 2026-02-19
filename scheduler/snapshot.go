package scheduler

import "time"

type Snapshot struct {
	GeneratedAt time.Time              `json:"generated_at"`
	StartedAt   time.Time              `json:"started_at"`
	Counts      SnapshotCounts         `json:"counts"`
	Checks      map[string]CheckStatus `json:"checks"`
}

type SnapshotCounts struct {
	Total   int `json:"total"`
	Running int `json:"running"`
	Queued  int `json:"queued"`
	Pass    int `json:"pass"`
	Fail    int `json:"fail"`
	Warn    int `json:"warn"`
}

type CheckStatus struct {
	Name           string              `json:"name"`
	NextRun        time.Time           `json:"next_run"`
	LastRun        *time.Time          `json:"last_run,omitempty"`
	LastStatus     string              `json:"last_status"`
	LastDurationMs int64               `json:"last_duration_ms"`
	Running        bool                `json:"running"`
	Queued         bool                `json:"queued"`
	ScheduleType   ScheduleType        `json:"schedule_type"`
	History        []CheckHistoryEntry `json:"history,omitempty"`
}

func (s *Scheduler) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := Snapshot{
		GeneratedAt: time.Now().In(s.location),
		StartedAt:   s.startTime,
		Counts: SnapshotCounts{
			Total:   len(s.checks),
			Running: s.runningChecks,
			Queued:  len(s.queue),
		},
		Checks: make(map[string]CheckStatus, len(s.checks)),
	}

	for name, check := range s.checks {
		lastRun := copyTimePtr(check.LastRun)
		history := copyHistory(check.History)
		snapshot.Checks[name] = CheckStatus{
			Name:           name,
			NextRun:        check.NextRun,
			LastRun:        lastRun,
			LastStatus:     check.LastStatus,
			LastDurationMs: check.LastDuration.Milliseconds(),
			Running:        check.Running,
			Queued:         check.IsQueued,
			ScheduleType:   check.ScheduleType,
			History:        history,
		}
		switch check.LastStatus {
		case "pass":
			snapshot.Counts.Pass++
		case "fail":
			snapshot.Counts.Fail++
		case "warn":
			snapshot.Counts.Warn++
		}
	}

	return snapshot
}

func copyTimePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}

func copyHistory(entries []CheckHistoryEntry) []CheckHistoryEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]CheckHistoryEntry, len(entries))
	copy(out, entries)
	return out
}
