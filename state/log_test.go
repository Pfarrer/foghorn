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
