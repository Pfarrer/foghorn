## Context

The daemon normally runs the scheduler's tick loop indefinitely. One-shot mode is a special execution path that loads config, iterates all enabled checks once (bypassing recurring schedule logic), evaluates results, logs outcomes, and exits.

## Goals / Non-Goals

**Goals:**
- Add a CLI flag to enable one-shot mode
- Execute each enabled check exactly once
- Use existing check execution, timeout, and evaluation paths
- Exit with code reflecting aggregate success/failure
- Preserve normal logging and state updates

**Non-Goals:**
- Replacing the scheduler for normal runs
- Adding special result aggregation beyond exit codes

## Decisions

**Flag vs. config**: CLI flag (`--one-shot`) rather than config field, as it's an execution mode rather than persistent behavior.

**Use existing executor**: Reuse `DockerExecutor.Execute()` directly rather than routing through the scheduler.

**Parallel execution**: Run checks in parallel up to `max_concurrent_checks` for faster execution.

**Exit code logic**: 0 if all checks pass/warn; 1 if any check fails/unknown; 2 for config errors or missing dependencies.

## Risks / Trade-offs

- [Parallel resource usage] → All checks running at once may exhaust Docker; mitigated by respecting `max_concurrent_checks`
