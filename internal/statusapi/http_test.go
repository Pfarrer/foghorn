package statusapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/scheduler"
)

func TestClientGetStatus(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	want := scheduler.Snapshot{
		GeneratedAt: now,
		StartedAt:   now.Add(-time.Minute),
		Counts: scheduler.SnapshotCounts{
			Total: 2,
			Pass:  1,
			Fail:  1,
		},
		Checks: map[string]scheduler.CheckStatus{
			"a": {
				Name:       "a",
				LastStatus: "pass",
			},
		},
	}
	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return want
	}))
	defer server.Close()

	client := NewClient(server.URL)
	got, err := client.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}
	if got.Counts.Total != want.Counts.Total {
		t.Fatalf("Counts.Total = %d, want %d", got.Counts.Total, want.Counts.Total)
	}
	if got.Checks["a"].LastStatus != "pass" {
		t.Fatalf("Checks[a].LastStatus = %q, want pass", got.Checks["a"].LastStatus)
	}
}

func TestStatusPathMethodNotAllowed(t *testing.T) {
	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}))
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+StatusPath, nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}
