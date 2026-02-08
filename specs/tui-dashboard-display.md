# TUI Dashboard Display

## Category
functional

## Description
Add a read-only terminal user interface (TUI) dashboard for real-time monitoring of check status and scheduler activity, including last run date/time, time until next run, and recent result history

## Usage Steps
1. Run Foghorn with `--tui` flag to enable TUI mode
2. Observe dashboard showing:
   - Summary counters (total checks, running, queued, last completed)
   - List of all checks with status indicators
   - Time remaining until next scheduled run for each check
   - Last check result (pass/fail/warn/error) per check
   - Last run date/time
   - Recent result history (last 10 runs)
3. Dashboard refreshes automatically every second
4. Press Ctrl+C to exit (no other interaction required)

## Implementation Notes
- Use Go TUI library (e.g., bubbletea, termui, or similar)
- TUI mode is mutually exclusive with CLI output mode
- Display sections:
  - Header: Foghorn title, uptime, configured log level
  - Summary bar: Total checks | Running | Queued | Pass | Fail | Warn
  - Check list: Table with columns:
    - Check name
    - Status (running icon, queued icon, idle)
    - Last result (✓ pass, ✗ fail, ⚠ warn, ? unknown)
    - Last run date/time
    - Next run time (countdown)
    - Result history (last 10 runs)
  - Footer: Refresh interval, help text (Ctrl+C to exit)
- Status indicators:
  - Running: ⟳ or ▶
  - Queued: ⏳
  - Idle: • or space
  - Pass: ✓
  - Fail: ✗
  - Warn: ⚠
  - Unknown: ?
- Refresh rate: 1 second (configurable)
- No user interaction - read-only display
- Must handle terminal resizing gracefully
- Colors for visual clarity (green=pass, red=fail, yellow=warn, etc.)
- Show max 20 checks at once with scrolling if more checks exist
- Time format: absolute date/time for last run (local time), countdown for next run (e.g., "in 45s")
- History format: last 10 results as compact symbols (e.g., "✓ ✓ ⚠ ✗ ?")

## Acceptance Criteria
- [x] `--tui` flag enables TUI dashboard mode
- [x] Header shows Foghorn title and uptime
- [x] Summary bar shows accurate counters (total, running, queued, pass, fail, warn)
- [x] Check list displays all configured checks
- [x] Check status indicator shows correct state (running, queued, idle)
- [x] Last check result shows correct status symbol (✓, ✗, ⚠, ?)
- [x] Last run time shows absolute date/time (local time)
- [x] Next run time shows countdown (e.g., "in 45s")
- [x] Check history shows last 10 results in order
- [x] Dashboard refreshes every second automatically
- [x] Dashboard handles terminal resize without crashing
- [x] No user interaction is required or possible (read-only)
- [x] Ctrl+C exits cleanly
- [x] Check list scrolls if more than 20 checks
- [x] Status colors are used appropriately (green=pass, red=fail, yellow=warn)
- [x] Works with both interval-based and cron-scheduled checks

Passes: true
