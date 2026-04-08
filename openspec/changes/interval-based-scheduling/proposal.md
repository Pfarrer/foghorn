## Why

Interval-based scheduling is implemented in `scheduler/scheduler.go` via `parseInterval` and `ScheduleTypeInterval` handling, with a legacy spec at `specs/interval-based-scheduling.md`. Migrating to OpenSpec formalizes this feature.

## What Changes

- Formalize the existing interval scheduling feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `interval-based-scheduling`: Time-duration-based scheduling as an alternative to cron expressions, supporting seconds, minutes, hours, and days

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/interval-based-scheduling.md` — legacy spec to be removed
- No code changes
