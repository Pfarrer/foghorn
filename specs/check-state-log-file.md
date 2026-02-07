# Check State Log File

## Category
config

## Description
Add a CLI argument to persist recent check results to a state log file and restore scheduling on startup.

## Usage Steps
1. Set a state log file path via a CLI argument.
2. Set the retention period for stored results.
3. Start Foghorn and run checks.
4. Restart Foghorn and observe restored intervals from the last results.

## Implementation Notes
- Add CLI argument for `state_log_file` (path).
- Add config field for `state_log_period` (duration).
- Persist check results for the most recent retention period only.
- Delete records older than the configured period on write and on startup.
- Read the state log file on startup and seed scheduler state from last results.
- Handle missing or corrupt state log file gracefully.

## Acceptance Criteria
- [ ] CLI supports `state_log_file`.
- [ ] Config supports `state_log_period`.
- [ ] Results are persisted to the state log file.
- [ ] Records older than the retention period are removed automatically.
- [ ] On restart, the scheduler restores intervals using the last stored results.
- [ ] Missing or invalid state log file does not crash startup.
