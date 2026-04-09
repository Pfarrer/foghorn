## Context

The state log (`state/log.go`) persists check results as JSON lines to a file. Each record contains `check_name`, `status`, `duration_ms`, and `completed_at`. Records older than the retention period are pruned on read and write. The daemon uses `LatestByCheck` to restore scheduling state from the most recent result per check. File locking via `syscall.Flock` prevents concurrent access.

## Goals / Non-Goals

**Goals:**
- Document the state log behavior as formal OpenSpec artifacts

**Non-Goals:**
- Adding new record fields, changing the format, or adding compaction strategies

## Decisions

**JSON lines format**: Each record is a single JSON object on one line, enabling append-only writes and line-by-line parsing.

**Full-file rewrite on append**: Rather than appending to the end, the file is rewritten after filtering expired records. This keeps the file size bounded.

**File locking**: Uses `syscall.Flock` with `LOCK_EX|LOCK_NB` to prevent multiple daemon instances from writing to the same file.

**Graceful degradation**: A missing or corrupt state log does not crash the daemon; it starts fresh.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
