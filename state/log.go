package state

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type Record struct {
	CheckName   string    `json:"check_name"`
	Status      string    `json:"status"`
	DurationMs  int64     `json:"duration_ms"`
	CompletedAt time.Time `json:"completed_at"`
}

type StateLog struct {
	path      string
	retention time.Duration
	lockFile  *os.File
	mu        sync.Mutex
}

func Open(path string, retention time.Duration) (*StateLog, error) {
	if path == "" {
		return nil, fmt.Errorf("state log path is required")
	}
	if retention <= 0 {
		return nil, fmt.Errorf("state log retention must be positive")
	}

	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create state log directory: %w", err)
		}
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open state log file: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("state log file is locked by another process")
	}

	return &StateLog{
		path:      path,
		retention: retention,
		lockFile:  file,
	}, nil
}

func (s *StateLog) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.lockFile == nil {
		return nil
	}
	_ = syscall.Flock(int(s.lockFile.Fd()), syscall.LOCK_UN)
	err := s.lockFile.Close()
	s.lockFile = nil
	return err
}

func (s *StateLog) RecordResult(checkName string, status string, duration time.Duration, completedAt time.Time) error {
	if completedAt.IsZero() {
		completedAt = time.Now().UTC()
	}
	record := Record{
		CheckName:   checkName,
		Status:      status,
		DurationMs:  duration.Milliseconds(),
		CompletedAt: completedAt.UTC(),
	}
	return s.Append(record)
}

func (s *StateLog) Load() ([]Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	records, err := s.readAll()
	if err != nil {
		return nil, err
	}

	filtered := s.filter(records, time.Now().UTC())
	if len(filtered) != len(records) {
		if err := s.writeAll(filtered); err != nil {
			return nil, err
		}
	}

	return filtered, nil
}

func (s *StateLog) Append(record Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	records, readErr := s.readAll()
	if readErr != nil {
		records = nil
	}

	filtered := s.filter(records, now)
	if record.CompletedAt.IsZero() {
		record.CompletedAt = now
	}
	if s.retention > 0 {
		cutoff := now.Add(-s.retention)
		if record.CompletedAt.Before(cutoff) {
			if err := s.writeAll(filtered); err != nil {
				return err
			}
			if readErr != nil {
				return fmt.Errorf("state log is corrupt: %w", readErr)
			}
			return nil
		}
	}

	filtered = append(filtered, record)
	if err := s.writeAll(filtered); err != nil {
		return err
	}

	if readErr != nil {
		return fmt.Errorf("state log is corrupt: %w", readErr)
	}

	return nil
}

func LatestByCheck(records []Record) map[string]Record {
	latest := make(map[string]Record, len(records))
	for _, record := range records {
		if record.CheckName == "" {
			continue
		}
		existing, ok := latest[record.CheckName]
		if !ok || record.CompletedAt.After(existing.CompletedAt) {
			latest[record.CheckName] = record
		}
	}
	return latest
}

func (s *StateLog) filter(records []Record, now time.Time) []Record {
	if s.retention <= 0 {
		return records
	}
	cutoff := now.Add(-s.retention)
	filtered := make([]Record, 0, len(records))
	for _, record := range records {
		if record.CompletedAt.Before(cutoff) {
			continue
		}
		filtered = append(filtered, record)
	}
	return filtered
}

func (s *StateLog) readAll() ([]Record, error) {
	if s.lockFile == nil {
		return nil, fmt.Errorf("state log is closed")
	}
	if _, err := s.lockFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(s.lockFile)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var records []Record
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var record Record
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (s *StateLog) writeAll(records []Record) error {
	if s.lockFile == nil {
		return fmt.Errorf("state log is closed")
	}
	if err := s.lockFile.Truncate(0); err != nil {
		return err
	}
	if _, err := s.lockFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	writer := bufio.NewWriter(s.lockFile)
	for _, record := range records {
		payload, err := json.Marshal(record)
		if err != nil {
			return err
		}
		if _, err := writer.Write(payload); err != nil {
			return err
		}
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}

	return s.lockFile.Sync()
}
