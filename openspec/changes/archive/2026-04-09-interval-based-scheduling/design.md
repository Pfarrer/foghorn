## Context

The scheduler supports two schedule types: `ScheduleTypeCron` and `ScheduleTypeInterval`. Interval checks are detected when a config implements `IntervalCheckConfig` and returns `ScheduleTypeInterval`. The `parseInterval` function converts strings like `5m`, `1h`, `30s`, `2d` into `time.Duration`.

## Goals / Non-Goals

**Goals:**
- Document the existing interval scheduling behavior as formal OpenSpec artifacts

**Non-Goals:**
- Adding new time units or composite durations (e.g., `1h30m`)

## Decisions

**Simple single-unit format**: Intervals use a single number and unit suffix (e.g., `5m`). No composite durations like Go's `1h30m` are supported.

**Immediate first run**: On empty state, interval checks run immediately. Subsequent runs are scheduled by adding the interval duration to the previous run time.

**Coexistence with cron**: Both cron and interval checks can exist in the same configuration and are managed by the same scheduler.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
