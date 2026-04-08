## Why

The CLI argument parsing is implemented in `internal/daemon/app.go` using Go's `flag` package with a legacy spec at `specs/basic-cli-arguments.md`. Migrating to OpenSpec formalizes the CLI interface for maintenance and consistency.

## What Changes

- Formalize the existing CLI arguments feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes to the daemon CLI

## Capabilities

### New Capabilities
- `basic-cli-arguments`: Command-line argument parsing for the foghorn-daemon binary, supporting help, config path, verbose mode, dry-run, log level, and related flags

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/basic-cli-arguments.md` — legacy spec to be removed
- No code changes
