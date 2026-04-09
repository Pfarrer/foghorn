## Why

The Docker check interface contract is implemented in `executor/docker.go` (env var injection and JSON output parsing) with a legacy spec at `specs/docker-check-interface.md`. Migrating to OpenSpec formalizes the contract.

## What Changes

- Formalize the existing check interface contract as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `docker-check-interface`: Defines the contract between Foghorn and check containers — environment variable inputs, JSON output format, exit code semantics, and error handling

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/docker-check-interface.md` — legacy spec to be removed
- No code changes
