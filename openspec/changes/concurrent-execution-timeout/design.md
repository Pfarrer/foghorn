## Context

The executor (`executor/docker.go`) runs check containers with `context.WithTimeout`. Per-check timeout is parsed from `checkConfig.Timeout`, falling back to a 30s default. The scheduler (`scheduler/scheduler.go`) enforces concurrency limits via a queue and `runningChecks` counter. Timeout events are logged and containers are terminated.

## Goals / Non-Goals

**Goals:**
- Document existing concurrency and timeout behavior as formal OpenSpec artifacts

**Non-Goals:**
- Modifying timeout logic, concurrency limits, or Docker container lifecycle

## Decisions

**Go context-based timeouts**: Uses `context.WithTimeout` to bound container execution. When the context expires, the container is removed.

**Queue-based concurrency**: When `maxConcurrentChecks` is reached, due checks enter a queue and are dispatched as slots free up.

**FOGHORN_TIMEOUT env var**: The parsed per-check timeout is injected as `FOGHORN_TIMEOUT` into the container's environment, allowing check containers to self-regulate.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
