## Context

The executor (`executor/docker.go`) captures container stdout/stderr and logs it based on the configured `check_container_debug_output` mode. Three modes are supported: `off` (default, no output), `on_failure` (log on check failure/timeout), and `always` (log every run). Output is redacted to remove secrets and auth patterns, then truncated to a configurable max character limit.

## Goals / Non-Goals

**Goals:**
- Document the debug output feature as formal OpenSpec artifacts

**Non-Goals:**
- Adding new modes or changing redaction patterns

## Decisions

**Three-mode model**: `off` is the default to avoid noisy logs. `on_failure` logs only on failures/timeouts. `always` logs everything.

**Per-check override**: Each check can override the global `check_container_debug_output` setting.

**Regex-based redaction**: Secret values and common auth patterns (Authorization headers, passwords, Bearer tokens) are replaced with `[REDACTED]` before logging.

**Truncation**: Output is truncated to a configurable max character limit (default 4096) to avoid oversized log lines.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
