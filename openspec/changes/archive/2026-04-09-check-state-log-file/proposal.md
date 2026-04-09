## Why

The check state log file is implemented in `state/log.go` with a legacy spec at `specs/check-state-log-file.md`. Migrating to OpenSpec formalizes the state persistence and restoration behavior.

## What Changes

- Formalize the existing state log feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `check-state-log-file`: Persistent state log that records check results to a JSON-lines file with configurable retention, and restores scheduling state on daemon restart

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/check-state-log-file.md` — legacy spec to be removed
- No code changes
