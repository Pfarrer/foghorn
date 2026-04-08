## Why

State file locking is implemented in `state/log.go` using `syscall.Flock` with a legacy spec at `specs/state-file-locking.md`. Migrating to OpenSpec formalizes the locking behavior.

## What Changes

- Formalize the state file locking feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `state-file-locking`: Exclusive file locking on the state log to prevent multiple daemon instances from sharing the same file

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/state-file-locking.md` — legacy spec to be removed
- No code changes
