## Why

The TUI dashboard is implemented in `tui/` and `cmd/foghorn-tui/` using Bubble Tea, with a legacy spec at `specs/tui-dashboard-display.md`. Migrating to OpenSpec formalizes the dashboard behavior.

## What Changes

- Formalize the existing TUI dashboard as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `tui-dashboard-display`: Read-only terminal dashboard for real-time monitoring of check status, scheduler activity, and result history via the daemon status API

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/tui-dashboard-display.md` — legacy spec to be removed
- No code changes
