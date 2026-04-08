## Context

The scheduler (`scheduler/scheduler.go`) manages check execution with concurrency limiting and priority ordering. The `max_concurrent_checks` config field (`config/types.go`) controls how many checks run simultaneously. When unset or 0, concurrency is unlimited.

## Goals / Non-Goals

**Goals:**
- Document the existing execution ordering and concurrency behavior as formal OpenSpec artifacts
- Preserve all current behavior in testable spec requirements

**Non-Goals:**
- Modifying scheduling or concurrency logic
- Adding new scheduling strategies or features

## Decisions

**Priority by interval length**: When multiple checks are due, longer-interval checks run first (higher `priorityDuration`). This ensures infrequent checks get priority over frequent ones.

**Stable tie-breaker**: When two checks have equal priority duration, alphabetical name ordering is used (`name < name`).

**Immediate run for interval checks**: On empty state (no prior execution history), interval-based checks run immediately without initial delay. Cron checks always wait for their scheduled time.

**Queue-based concurrency**: When `max_concurrent_checks` > 0, due checks enter a queue. The `processQueue` method dequeues and dispatches checks up to the concurrency limit.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
