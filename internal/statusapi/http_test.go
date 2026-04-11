package statusapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pfarrer/foghorn/scheduler"
	"github.com/pfarrer/foghorn/state"
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
	}, nil))
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

func TestStatusAPIResponseContainsNoSecrets(t *testing.T) {
	secretValue := "super-secret-api-key-99999"
	snapshot := scheduler.Snapshot{
		GeneratedAt: time.Now().UTC().Truncate(time.Second),
		Counts: scheduler.SnapshotCounts{
			Total: 1,
			Pass:  1,
		},
		Checks: map[string]scheduler.CheckStatus{
			"db-check": {
				Name:       "db-check",
				LastStatus: "pass",
			},
		},
	}

	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return snapshot
	}, nil))
	defer server.Close()

	client := NewClient(server.URL)
	got, err := client.GetStatus(context.Background())
	if err != nil {
		t.Fatalf("GetStatus() error = %v", err)
	}

	body, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	if strings.Contains(string(body), secretValue) {
		t.Fatalf("status API response should not contain secret: %s", string(body))
	}
}

func TestStatusPathMethodNotAllowed(t *testing.T) {
	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}, nil))
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

type mockHistoryQuerier struct {
	records []state.Record
	err     error
}

func (m *mockHistoryQuerier) Query(checkName string, start, end time.Time) ([]state.Record, error) {
	return m.records, m.err
}

func TestHistoryEndpointReturnsRecords(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	records := []state.Record{
		{CheckName: "test-check", Status: "pass", DurationMs: 100, CompletedAt: now},
	}
	querier := &mockHistoryQuerier{records: records}

	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}, querier))
	defer server.Close()

	resp, err := http.Get(server.URL + HistoryPath)
	if err != nil {
		t.Fatalf("GET history: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got []state.Record
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 || got[0].CheckName != "test-check" {
		t.Fatalf("expected 1 record for test-check, got %v", got)
	}
}

func TestHistoryEndpointWithCheckFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	records := []state.Record{
		{CheckName: "tls-check", Status: "pass", CompletedAt: now},
	}
	querier := &mockHistoryQuerier{records: records}

	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}, querier))
	defer server.Close()

	resp, err := http.Get(server.URL + HistoryPath + "?check=tls-check")
	if err != nil {
		t.Fatalf("GET history: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHistoryEndpointWithTimeRange(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	querier := &mockHistoryQuerier{records: nil}

	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}, querier))
	defer server.Close()

	start := now.Add(-1 * time.Hour).Format(time.RFC3339)
	end := now.Format(time.RFC3339)
	resp, err := http.Get(server.URL + HistoryPath + "?start=" + start + "&end=" + end)
	if err != nil {
		t.Fatalf("GET history: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHistoryEndpointNotAvailable(t *testing.T) {
	server := httptest.NewServer(NewHandler(func() scheduler.Snapshot {
		return scheduler.Snapshot{}
	}, nil))
	defer server.Close()

	resp, err := http.Get(server.URL + HistoryPath)
	if err != nil {
		t.Fatalf("GET history: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusServiceUnavailable)
	}
}
