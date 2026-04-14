## Context

The check interface defines four status values: `pass`, `fail`, `warn`, `unknown`. Check containers output one of these in their JSON response. Currently there is no way to express "the check could not complete due to a technical problem in the check itself" (e.g., network error fetching a page, browser crash, container misconfiguration). These are reported as `fail`, conflating a broken check with a broken property.

The status is stored as a plain string throughout the codebase — `state.Record.Status`, `scheduler.CheckState.LastStatus`, `executor.CheckResult.Status`, `snapshot.CheckStatus.LastStatus` — with `switch` statements in `scheduler.go`, `snapshot.go`, `daemon/app.go`, and `tui/model.go`.

## Goals / Non-Goals

**Goals:**
- Add `error` as a first-class status value meaning "the check could not finish due to a technical error in the container"
- Update all status-handling code to recognize `error`
- Update snapshot counts to include `error`
- Update TUI to render `error` distinctly from `fail`

**Non-Goals:**
- Changing semantics of existing statuses (`pass`, `fail`, `warn`, `unknown`)
- Requiring existing check containers to adopt `error` (backward compatible)
- Adding sub-categorization of errors (containers can use `data.error_type` freely)

## Decisions

### 1. `error` status for technical failures

**Decision**: Add `"error"` as a new status value in the check contract.

**Rationale**: `fail` means "the checked property has a problem." `error` means "the check itself could not run." This distinction is critical for monitoring: a network outage should not look identical to a detected page change. Downstream consumers (alerts, dashboards) need to react differently.

**Alternatives**:
- **Reuse `unknown`**: `unknown` already exists but is used as a fallback/default for unstarted checks and missing data. Overloading it with "technical error" loses the fallback semantics.
- **Use `data.error_type` only**: Requires parsing the `data` object to distinguish error categories. A top-level status is simpler and consistent with existing contract.

### 2. `error` counts as unhealthy in one-shot mode

**Decision**: In one-shot mode, `error` status SHALL cause a non-zero exit code, same as `fail` and `unknown`.

**Rationale**: If a check couldn't run, the monitoring result is not healthy. Silent failures are worse than false alerts.

### 3. TUI rendering

**Decision**: Render `error` with a distinct symbol and color (`!` in yellow/amber), different from `fail` (`✗` in red).

**Rationale**: Users need to visually distinguish "check failed" from "check errored" at a glance.

## Risks / Trade-offs

- **Breaking change for external consumers** → Any external tool parsing Foghorn status JSON must handle the new `error` value. Mitigation: unknown status values already fall through to `default` in all `switch` statements, so existing code won't break — it just won't specially handle `error`.
- **Existing containers unaffected** → Containers that don't emit `error` continue working identically. Only new containers (e.g., `web-change-check`) will use it.
