## ADDED Requirements

### Requirement: Extended history records
Each execution record SHALL include: check name, status, start time, completion time, duration in milliseconds, and error details (when applicable). Records SHALL be persisted to the state log file.

#### Scenario: Successful run recorded
- **WHEN** a check completes with status `pass` in 250ms
- **THEN** a record is stored with check name, status `pass`, start time, completion time, duration 250ms, and empty error details

#### Scenario: Failed run with error details
- **WHEN** a check fails with an error message
- **THEN** the record includes the error details string

### Requirement: Retention policy
The system SHALL enforce a retention policy on execution history. Records older than the retention period SHALL be pruned on write. The retention period SHALL be configurable via `state_log_period`.

#### Scenario: Old records pruned
- **WHEN** the retention period is 24h and a record is 25h old
- **THEN** the record is removed during the next write

#### Scenario: Recent records retained
- **WHEN** a record is within the retention period
- **THEN** it is kept

### Requirement: Query by check name
The system SHALL support querying execution history filtered by check name.

#### Scenario: Filter by check name
- **WHEN** history is queried for check `tls-check`
- **THEN** only records for `tls-check` are returned

### Requirement: Query by time range
The system SHALL support querying execution history filtered by a time range (start time, end time).

#### Scenario: Filter by time range
- **WHEN** history is queried for the last hour
- **THEN** only records with completion time within the last hour are returned

#### Scenario: Combined filters
- **WHEN** history is queried for check `disk-check` within the last 24h
- **THEN** only matching records for that check within that time range are returned

### Requirement: History persistence across restarts
Execution history SHALL survive daemon restarts. On startup, the state log SHALL be loaded and history SHALL be available for querying.

#### Scenario: History survives restart
- **WHEN** the daemon is restarted
- **THEN** previously recorded execution history is available via queries
