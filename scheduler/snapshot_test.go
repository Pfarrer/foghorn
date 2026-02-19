package scheduler

import (
	"testing"
	"time"
)

func TestSnapshotIncludesCountsAndChecks(t *testing.T) {
	executor := &MockExecutor{}
	s := NewScheduler(executor, time.UTC, 0)

	check := &MockCheckConfig{
		name:     "snap-check",
		schedule: "*/5 * * * *",
		enabled:  true,
	}
	if err := s.AddCheck(check); err != nil {
		t.Fatalf("AddCheck() error = %v", err)
	}

	lastRun := time.Now().Add(-time.Minute).UTC()
	s.ApplyState(map[string]CheckState{
		"snap-check": {
			LastStatus: "pass",
			LastRun:    lastRun,
			History: []CheckHistoryEntry{
				{Status: "fail", CompletedAt: lastRun.Add(-time.Minute)},
				{Status: "pass", CompletedAt: lastRun},
			},
		},
	})

	snap := s.Snapshot()
	if snap.Counts.Total != 1 {
		t.Fatalf("Counts.Total = %d, want 1", snap.Counts.Total)
	}
	if snap.Counts.Pass != 1 {
		t.Fatalf("Counts.Pass = %d, want 1", snap.Counts.Pass)
	}

	checkSnap, ok := snap.Checks["snap-check"]
	if !ok {
		t.Fatalf("snapshot missing check snap-check")
	}
	if checkSnap.LastStatus != "pass" {
		t.Fatalf("LastStatus = %q, want pass", checkSnap.LastStatus)
	}
	if checkSnap.LastRun == nil {
		t.Fatalf("LastRun should not be nil")
	}
	if len(checkSnap.History) != 2 {
		t.Fatalf("History len = %d, want 2", len(checkSnap.History))
	}
}
