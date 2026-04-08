## Context

The scheduler (`scheduler/scheduler.go`) runs checks on a ticker-based loop. Each tick, it determines which checks are due by comparing current time against `NextRun`. Cron-scheduled checks use a custom cron parser to calculate next run times. The scheduler operates in a configured time location.

## Goals / Non-Goals

**Goals:**
- Document existing cron scheduling behavior as formal OpenSpec artifacts

**Non-Goals:**
- Modifying scheduling logic or adding new schedule types

## Decisions

**Custom cron parser**: Uses a local `ParseCronExpression` implementation supporting standard 5-field cron (minute, hour, day, month, day of week). No external cron library dependency.

**Ticker-based execution**: Uses `time.Ticker` with configurable interval (default 1s in the daemon). Each tick evaluates due checks.

**Time zone support**: Scheduler accepts a `*time.Location`. All time comparisons use `time.Now().In(location)`.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
