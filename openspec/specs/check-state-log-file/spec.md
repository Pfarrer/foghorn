## ADDED Requirements

### Requirement: State log configuration
The daemon SHALL accept `state_log_file` (file path) via CLI flag or config. When set, `state_log_period` (Go duration string, required when file is set) defines the retention period.

#### Scenario: Config validated
- **WHEN** `state_log_file` is set and `state_log_period` is a valid positive duration
- **THEN** the state log is opened and used for persistence

#### Scenario: Missing retention period
- **WHEN** `state_log_file` is set but `state_log_period` is missing
- **THEN** an error is returned indicating the period is required

#### Scenario: Invalid retention period
- **WHEN** `state_log_period` is not a valid positive duration
- **THEN** an error is returned

### Requirement: Result persistence
Each check result SHALL be recorded as a JSON line containing `check_name`, `status`, `duration_ms`, and `completed_at`. Records SHALL be written to the state log file.

#### Scenario: Result recorded
- **WHEN** a check completes with status `pass` and duration 250ms
- **THEN** a JSON record is appended to the state log file with those values

### Requirement: Retention pruning
Records older than the configured retention period SHALL be removed automatically on read and write.

#### Scenario: Old records pruned on load
- **WHEN** the state log contains records older than the retention period
- **THEN** `Load()` returns only records within the retention period

#### Scenario: Old records pruned on append
- **WHEN** a new result is appended and expired records exist
- **THEN** the file is rewritten with only non-expired records

### Requirement: State restoration on restart
On startup, the daemon SHALL read the state log, extract the latest result per check via `LatestByCheck`, and seed the scheduler with that state.

#### Scenario: Scheduler restored
- **WHEN** the daemon restarts and the state log has recent results
- **THEN** interval-based checks schedule from the last run time, not immediately

### Requirement: Graceful handling of missing/corrupt state log
A missing state log file SHALL not prevent startup. A corrupt state log SHALL be handled gracefully without crashing.

#### Scenario: Missing file
- **WHEN** `state_log_file` points to a non-existent path
- **THEN** the file is created and the daemon starts fresh

#### Scenario: Corrupt file
- **WHEN** the state log contains invalid JSON lines
- **THEN** the daemon logs a warning and starts fresh without crashing

### Requirement: File locking
The state log SHALL use file locking to prevent concurrent access by multiple daemon instances.

#### Scenario: Lock acquired
- **WHEN** the daemon opens the state log
- **THEN** an exclusive lock is acquired on the file

#### Scenario: Lock contention
- **WHEN** another process holds the lock
- **THEN** an error is returned indicating the file is locked
