## Why

Check container debug output is implemented in `executor/docker.go` with a legacy spec at `specs/check-container-debug-output.md`. Migrating to OpenSpec formalizes the debug logging and redaction behavior.

## What Changes

- Formalize the check container debug output feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `check-container-debug-output`: Configurable debug output modes (`off`, `on_failure`, `always`) for check containers with secret redaction and output truncation

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/check-container-debug-output.md` — legacy spec to be removed
- No code changes
