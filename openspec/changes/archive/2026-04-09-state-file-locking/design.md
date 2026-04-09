## Context

The state log (`state/log.go`) acquires an exclusive lock via `syscall.Flock(LOCK_EX|LOCK_NB)` when opened. The lock is held for the process lifetime and released on `Close()`.

## Goals / Non-Goals

**Goals:**
- Document the file locking behavior as formal OpenSpec artifacts

**Non-Goals:**
- Changing locking mechanism or adding distributed locking

## Decisions

**Non-blocking exclusive lock**: Uses `LOCK_NB` so the daemon fails immediately if another instance holds the lock, rather than blocking.

**Lock tied to file lifetime**: The lock is acquired on `Open()` and released on `Close()`, ensuring the daemon holds the lock for its entire lifetime.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
