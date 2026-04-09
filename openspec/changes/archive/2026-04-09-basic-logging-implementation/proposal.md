## Why

The logging system is implemented in `logger/logger.go` with a legacy spec at `specs/basic-logging-implementation.md`. Migrating to OpenSpec formalizes the logging behavior.

## What Changes

- Formalize the existing logging system as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `basic-logging-implementation`: Structured logging with configurable levels, verbose timestamps, and global logger

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/basic-logging-implementation.md` — legacy spec to be removed
- No code changes
