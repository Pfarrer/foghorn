## Why

Concurrent execution and timeout handling is implemented in the executor (`executor/docker.go`) and scheduler (`scheduler/scheduler.go`) with a legacy spec at `specs/concurrent-execution-timeout.md`. Migrating to OpenSpec formalizes this behavior.

## What Changes

- Formalize the existing concurrent execution and timeout feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `concurrent-execution-timeout`: Manages concurrent Docker container check executions with configurable timeouts, queuing when limits are reached, and terminating stuck checks

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/concurrent-execution-timeout.md` — legacy spec to be removed
- No code changes
