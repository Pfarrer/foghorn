## Why

The check execution ordering and concurrency feature is implemented in the scheduler (`scheduler/scheduler.go`) and config (`config/types.go`) with a legacy spec at `specs/check-execution-order-concurrency.md`. Migrating to OpenSpec formalizes the behavior for maintenance and consistency.

## What Changes

- Formalize the existing execution ordering and concurrency feature as an OpenSpec change with full artifacts
- Port the legacy spec into the OpenSpec spec format under `openspec/specs/check-execution-order-concurrency/`
- No functional changes to the scheduler or config

## Capabilities

### New Capabilities
- `check-execution-order-concurrency`: Limits concurrent check container execution and defines deterministic priority ordering when multiple checks are due simultaneously

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/check-execution-order-concurrency.md` — legacy spec to be removed
- No code changes to `scheduler/`, `config/`, or `internal/daemon/`
