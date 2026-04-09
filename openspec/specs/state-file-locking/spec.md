## ADDED Requirements

### Requirement: Exclusive lock on startup
The daemon SHALL acquire an exclusive lock on the state log file when opening it. The lock SHALL be held for the lifetime of the process.

#### Scenario: Lock acquired
- **WHEN** the daemon opens a state log file
- **THEN** an exclusive lock is acquired on the file

#### Scenario: Lock held during runtime
- **WHEN** the daemon is running with a state log
- **THEN** the lock remains held for the entire process lifetime

### Requirement: Second instance rejected
If a second daemon instance attempts to open the same state log file, it SHALL fail to acquire the lock and exit with a clear error message.

#### Scenario: Lock contention
- **WHEN** a second instance tries to open an already-locked state log file
- **THEN** an error is returned indicating the file is locked by another process

#### Scenario: Clear error message
- **WHEN** lock acquisition fails
- **THEN** the error message indicates the state log file is locked

### Requirement: Lock released on shutdown
When the daemon shuts down and closes the state log, the file lock SHALL be released.

#### Scenario: Lock released on close
- **WHEN** `StateLog.Close()` is called
- **THEN** `LOCK_UN` is called and the file descriptor is closed
