# State File Locking

## Category
config

## Description
Lock the state file during runtime and prevent multiple instances from sharing the same state file.

## Usage Steps
1. Set a state log file path in config.
2. Start Foghorn and keep it running.
3. Attempt to start a second instance using the same state file.
4. The second instance fails with a clear error message.

## Implementation Notes
- Acquire an exclusive lock on the state log file at startup.
- Hold the lock for the lifetime of the process.
- If the lock cannot be acquired, exit with a clear error message.
- Release the lock on shutdown.

## Acceptance Criteria
- [x] The state log file is locked while Foghorn is running.
- [x] A second instance using the same state file fails to start.
- [x] The failure reports a clear error message.
- [x] The lock is released on shutdown.

Passes: true
