## Why

The daemon/client modularization is implemented with `cmd/foghorn-daemon`, `cmd/foghorn-tui`, and `internal/` packages, with a legacy spec at `specs/daemon-client-modularization.md`. Migrating to OpenSpec formalizes the architecture.

## What Changes

- Formalize the daemon/client architecture as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `daemon-client-modularization`: Separation of Foghorn into a daemon process (scheduler/executor) and a standalone TUI client process communicating via a local status API

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/daemon-client-modularization.md` — legacy spec to be removed
- No code changes
