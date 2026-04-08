## Why

One-shot mode is spec'd but not implemented. It allows running all configured checks once without starting the scheduler, then exiting. This is useful for CI/CD pipelines, health checks triggered by external systems, and debugging. Creating this as an OpenSpec change defines the requirements and implementation plan.

## What Changes

- Add a CLI flag to enable one-shot mode
- Bypass the scheduler's recurring tick loop
- Execute each enabled check exactly once
- Use existing executor, timeout, and evaluation logic
- Exit with code 0 on all passes, non-zero on any failure

## Capabilities

### New Capabilities
- `one-shot-mode`: Run all configured checks once and exit, with exit code reflecting aggregate result

## Impact

- `internal/daemon/app.go` — add CLI flag and one-shot execution flow
- `specs/one-shot-mode.md` — legacy spec to be removed
