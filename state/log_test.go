package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStateLogRetention(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("open state log: %v", err)
	}
	defer log.Close()

	now := time.Now().UTC()
	oldRecord := Record{CheckName: "old", Status: "pass", DurationMs: 100, CompletedAt: now.Add(-2 * time.Hour)}
	newRecord := Record{CheckName: "new", Status: "fail", DurationMs: 200, CompletedAt: now.Add(-10 * time.Minute)}

	if err := log.Append(oldRecord); err != nil {
		t.Fatalf("append old record: %v", err)
	}
	if err := log.Append(newRecord); err != nil {
		t.Fatalf("append new record: %v", err)
	}

	records, err := log.Load()
	if err != nil {
		t.Fatalf("load state log: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].CheckName != "new" {
		t.Fatalf("expected new record, got %s", records[0].CheckName)
	}
}

func TestLatestByCheck(t *testing.T) {
	now := time.Now().UTC()
	records := []Record{
		{CheckName: "a", Status: "pass", CompletedAt: now.Add(-2 * time.Minute)},
		{CheckName: "a", Status: "fail", CompletedAt: now.Add(-1 * time.Minute)},
		{CheckName: "b", Status: "warn", CompletedAt: now.Add(-3 * time.Minute)},
	}

	latest := LatestByCheck(records)
	if latest["a"].Status != "fail" {
		t.Fatalf("expected latest status fail, got %s", latest["a"].Status)
	}
	if latest["b"].Status != "warn" {
		t.Fatalf("expected latest status warn, got %s", latest["b"].Status)
	}
}

func TestExclusiveLockAcquiredOnOpen(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")

	log, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("open state log: %v", err)
	}
	defer log.Close()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open file for lock check: %v", err)
	}
	defer file.Close()

	if log.lockFile == nil {
		t.Fatal("expected lockFile to be set after Open")
	}
}

func TestSecondInstanceRejected(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")

	log1, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("open first state log: %v", err)
	}
	defer log1.Close()

	_, err = Open(path, time.Hour)
	if err == nil {
		t.Fatal("expected error when opening locked state log from second instance")
	}
}

func TestLockReleasedOnClose(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")

	log1, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("open state log: %v", err)
	}

	if err := log1.Close(); err != nil {
		t.Fatalf("close state log: %v", err)
	}

	log2, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("expected to re-open state log after close: %v", err)
	}
	defer log2.Close()
}

func TestLoadCorruptStateLog(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	if err := os.WriteFile(path, []byte("not-json\n"), 0o644); err != nil {
		t.Fatalf("write corrupt file: %v", err)
	}

	log, err := Open(path, time.Hour)
	if err != nil {
		t.Fatalf("open state log: %v", err)
	}
	defer log.Close()

	if _, err := log.Load(); err == nil {
		t.Fatalf("expected load error for corrupt log")
	}
}

func TestExtendedRecordPersistenceAndLoading(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	startTime := time.Now().UTC().Add(-250 * time.Millisecond).Truncate(time.Millisecond)
	completedAt := time.Now().UTC().Truncate(time.Millisecond)

	err = log.RecordResultFull("test-check", "pass", 250*time.Millisecond, completedAt, startTime, "")
	if err != nil {
		t.Fatalf("RecordResultFull: %v", err)
	}

	err = log.RecordResultFull("fail-check", "fail", 100*time.Millisecond, completedAt, time.Time{}, "container exited with code 1")
	if err != nil {
		t.Fatalf("RecordResultFull: %v", err)
	}

	records, err := log.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	var passRec, failRec *Record
	for i := range records {
		if records[i].CheckName == "test-check" {
			passRec = &records[i]
		} else if records[i].CheckName == "fail-check" {
			failRec = &records[i]
		}
	}

	if passRec == nil {
		t.Fatal("pass record not found")
	}
	if !passRec.StartedAt.Equal(startTime) {
		t.Errorf("started_at = %v, want %v", passRec.StartedAt, startTime)
	}
	if passRec.ErrorDetails != "" {
		t.Errorf("error_details should be empty for pass, got %q", passRec.ErrorDetails)
	}

	if failRec == nil {
		t.Fatal("fail record not found")
	}
	if failRec.ErrorDetails != "container exited with code 1" {
		t.Errorf("error_details = %q, want %q", failRec.ErrorDetails, "container exited with code 1")
	}
}

func TestBackwardCompatibilityOldRecordsLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")

	oldJSON := `{"check_name":"old-check","status":"pass","duration_ms":100,"completed_at":"2026-01-01T00:00:00Z"}` + "\n"
	if err := os.WriteFile(path, []byte(oldJSON), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	log, err := Open(path, 365*24*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	records, err := log.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].CheckName != "old-check" {
		t.Errorf("check_name = %q, want %q", records[0].CheckName, "old-check")
	}
	if !records[0].StartedAt.IsZero() {
		t.Errorf("started_at should be zero for old record, got %v", records[0].StartedAt)
	}
	if records[0].ErrorDetails != "" {
		t.Errorf("error_details should be empty for old record, got %q", records[0].ErrorDetails)
	}
}

func TestQueryByCheckName(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	now := time.Now().UTC()
	log.Append(Record{CheckName: "tls-check", Status: "pass", CompletedAt: now.Add(-5 * time.Minute)})
	log.Append(Record{CheckName: "disk-check", Status: "fail", CompletedAt: now.Add(-3 * time.Minute)})
	log.Append(Record{CheckName: "tls-check", Status: "warn", CompletedAt: now.Add(-1 * time.Minute)})

	records, err := log.Query("tls-check", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records for tls-check, got %d", len(records))
	}
	for _, r := range records {
		if r.CheckName != "tls-check" {
			t.Errorf("got check %q, want tls-check", r.CheckName)
		}
	}
}

func TestQueryByTimeRange(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	now := time.Now().UTC()
	log.Append(Record{CheckName: "a", Status: "pass", CompletedAt: now.Add(-2 * time.Hour)})
	log.Append(Record{CheckName: "b", Status: "pass", CompletedAt: now.Add(-30 * time.Minute)})
	log.Append(Record{CheckName: "c", Status: "pass", CompletedAt: now.Add(-5 * time.Minute)})

	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Minute)
	records, err := log.Query("", start, end)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records in time range, got %d", len(records))
	}
	for _, r := range records {
		if r.CheckName == "a" {
			t.Errorf("record 'a' should be outside time range")
		}
	}
}

func TestQueryCombinedFilters(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, 24*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	now := time.Now().UTC()
	log.Append(Record{CheckName: "disk-check", Status: "pass", CompletedAt: now.Add(-5 * time.Hour)})
	log.Append(Record{CheckName: "disk-check", Status: "fail", CompletedAt: now.Add(-30 * time.Minute)})
	log.Append(Record{CheckName: "tls-check", Status: "pass", CompletedAt: now.Add(-10 * time.Minute)})

	start := now.Add(-1 * time.Hour)
	records, err := log.Query("disk-check", start, time.Time{})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record for disk-check in last hour, got %d", len(records))
	}
	if records[0].Status != "fail" {
		t.Errorf("status = %q, want fail", records[0].Status)
	}
}

func TestRetentionWithExtendedRecords(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "state.log")
	log, err := Open(path, 1*time.Hour)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer log.Close()

	now := time.Now().UTC()
	oldRecord := Record{
		CheckName:    "old",
		Status:       "pass",
		DurationMs:   100,
		CompletedAt:  now.Add(-2 * time.Hour),
		StartedAt:    now.Add(-2 * time.Hour).Add(-100 * time.Millisecond),
		ErrorDetails: "old error",
	}
	newRecord := Record{
		CheckName:    "new",
		Status:       "fail",
		DurationMs:   200,
		CompletedAt:  now.Add(-10 * time.Minute),
		StartedAt:    now.Add(-10 * time.Minute).Add(-200 * time.Millisecond),
		ErrorDetails: "new error",
	}

	if err := log.Append(oldRecord); err != nil {
		t.Fatalf("append old: %v", err)
	}
	if err := log.Append(newRecord); err != nil {
		t.Fatalf("append new: %v", err)
	}

	records, err := log.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record after retention, got %d", len(records))
	}
	if records[0].CheckName != "new" {
		t.Errorf("expected new record, got %s", records[0].CheckName)
	}
	if records[0].ErrorDetails != "new error" {
		t.Errorf("error_details = %q, want %q", records[0].ErrorDetails, "new error")
	}
	if records[0].StartedAt.IsZero() {
		t.Error("started_at should not be zero")
	}
}
