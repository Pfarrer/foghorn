## Why

The scheduler with cron support is implemented in `scheduler/scheduler.go` with a legacy spec at `specs/basic-scheduler-cron.md`. Migrating to OpenSpec formalizes the scheduling behavior.

## What Changes

- Formalize the existing cron scheduler feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `basic-scheduler-cron`: Core scheduler that triggers check executions based on cron expressions, calculating next run times and respecting time zones

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/basic-scheduler-cron.md` — legacy spec to be removed
- No code changes
