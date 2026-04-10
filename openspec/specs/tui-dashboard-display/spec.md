# tui-dashboard-display Specification

## Purpose
TBD - created by archiving change tui-dashboard-display. Update Purpose after archive.
## Requirements
### Requirement: Dashboard startup
The `foghorn-tui` binary SHALL start a read-only TUI dashboard connecting to the daemon status API. It SHALL accept `-u`/`--status-url` for the API endpoint and `-h`/`--help` for usage.

#### Scenario: Successful startup
- **WHEN** the daemon is running and `foghorn-tui` is started
- **THEN** the dashboard displays and refreshes automatically

#### Scenario: Daemon unreachable
- **WHEN** the daemon status API is not available
- **THEN** an error message is printed and the TUI exits with code 1

### Requirement: Header display
The dashboard header SHALL show the Foghorn title, uptime, and configured log level.

#### Scenario: Header rendered
- **WHEN** the dashboard is displayed
- **THEN** the header contains the title, uptime, and log level

### Requirement: Summary bar
The dashboard SHALL show a summary bar with counters: total checks, running, queued, pass, fail, warn.

#### Scenario: Summary counters accurate
- **WHEN** 10 checks exist, 1 is running, 0 queued, 7 pass, 2 fail
- **THEN** the summary bar reflects these counts

### Requirement: Check list
The dashboard SHALL display a list of all checks with columns: check name, status indicator (running/queued/idle), last result (pass/fail/warn/unknown), last run date/time, next run countdown, and last 10 results as history symbols.

#### Scenario: Check list displays all checks
- **WHEN** the dashboard is loaded
- **THEN** all configured checks appear in the list

#### Scenario: Status indicators
- **WHEN** a check is running
- **THEN** a running indicator is displayed next to its name

#### Scenario: History symbols
- **WHEN** a check has results pass, pass, fail, warn, pass
- **THEN** the history column shows those results as compact symbols

### Requirement: Automatic refresh
The dashboard SHALL refresh automatically every second without user interaction.

#### Scenario: Auto-refresh
- **WHEN** a check status changes on the daemon
- **THEN** the dashboard reflects the change within approximately 1 second

### Requirement: Terminal resize handling
The dashboard SHALL handle terminal resizing without crashing.

#### Scenario: Terminal resized
- **WHEN** the terminal window is resized
- **THEN** the dashboard re-renders to fit the new dimensions

### Requirement: Scroll for many checks
When more than a screen-full of checks exist, the dashboard SHALL support scrolling.

#### Scenario: Many checks
- **WHEN** more checks exist than fit on screen
- **THEN** the user can scroll to see all checks

### Requirement: Clean exit
Pressing Ctrl+C SHALL exit the TUI cleanly.

#### Scenario: Ctrl+C exit
- **WHEN** the user presses Ctrl+C
- **THEN** the TUI exits cleanly without errors

### Requirement: Status colors
The dashboard SHALL use colors for visual clarity: green for pass, red for fail, yellow for warn.

#### Scenario: Color coding
- **WHEN** checks have pass, fail, and warn statuses
- **THEN** they are displayed in green, red, and yellow respectively

