# Check Execution Order and Concurrency

## Category
functional

## Description
Define execution concurrency limits and deterministic ordering for scheduled checks.

## Usage Steps
1. Set `max_concurrent_checks` in the config.
2. Start Foghorn.
3. Observe that checks are queued and executed in priority order.

## Implementation Notes
- Add a config property `max_concurrent_checks` (int, optional).
- Enforce at most `max_concurrent_checks` containers running at once.
- When multiple checks are due, prioritize the check with the longer interval (next planned run farther in the future).
- For equal intervals, keep a stable tie-breaker (e.g., name).
- On empty state, interval-based checks run immediately (no initial delay).
- Cron-scheduled checks are not affected by the immediate-run rule.

## Acceptance Criteria
- [x] `max_concurrent_checks` limits concurrent check containers.
- [x] Due checks execute in priority order (longer interval first).
- [x] Interval checks run immediately on empty state.
- [x] Cron checks still wait for their scheduled time.

## Passes
true
