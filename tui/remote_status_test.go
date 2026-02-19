package tui

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/internal/statusapi"
	"github.com/pfarrer/foghorn/scheduler"
)

func TestRemoteStatusReaderRefresh(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	server := httptest.NewServer(statusapi.NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{
			GeneratedAt: now,
			StartedAt:   now.Add(-2 * time.Minute),
			Counts: scheduler.SnapshotCounts{
				Total: 1,
				Pass:  1,
			},
			Checks: map[string]scheduler.CheckStatus{
				"api-check": {
					Name:       "api-check",
					LastStatus: "pass",
					History: []scheduler.CheckHistoryEntry{
						{Status: "pass", CompletedAt: now.Add(-time.Minute)},
					},
				},
			},
		}
	}))
	defer server.Close()

	reader, err := newRemoteStatusReader(server.URL)
	if err != nil {
		t.Fatalf("newRemoteStatusReader() error = %v", err)
	}
	if err := reader.Refresh(context.Background()); err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}

	total, _, _, pass, _, _ := reader.GetCounts()
	if total != 1 || pass != 1 {
		t.Fatalf("counts = total:%d pass:%d, want total:1 pass:1", total, pass)
	}
	checks := reader.GetAllChecks()
	if checks["api-check"] == nil {
		t.Fatalf("expected api-check in checks map")
	}
	if checks["api-check"].LastStatus != "pass" {
		t.Fatalf("LastStatus = %q, want pass", checks["api-check"].LastStatus)
	}
}
