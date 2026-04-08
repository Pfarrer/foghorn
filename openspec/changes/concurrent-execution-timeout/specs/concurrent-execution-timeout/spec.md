## ADDED Requirements

### Requirement: Per-check timeout configuration
Each check SHALL support an optional `timeout` field specifying a Go duration string. When set, the container execution SHALL be bounded by that duration. When unset, a 30-second default SHALL apply. The parsed timeout SHALL be injected as `FOGHORN_TIMEOUT` into the container environment.

#### Scenario: Per-check timeout set
- **WHEN** a check has `timeout: "60s"`
- **THEN** the container is executed with a 60-second context timeout and `FOGHORN_TIMEOUT=60s` in its environment

#### Scenario: Per-check timeout unset
- **WHEN** a check has no timeout configured
- **THEN** the container is executed with the 30-second default timeout

#### Scenario: Invalid timeout value
- **WHEN** a check has an unparseable timeout value
- **THEN** the executor falls back to the default timeout

### Requirement: Timeout container termination
When a container execution exceeds its timeout, the executor SHALL terminate the container and return a timeout error.

#### Scenario: Container exceeds timeout
- **WHEN** a check container runs longer than its configured timeout
- **THEN** the container is stopped/removed and a timeout error is returned

#### Scenario: Timeout logged
- **WHEN** a check times out
- **THEN** a warning is logged indicating the check name and timeout duration

### Requirement: Concurrency queuing
When the concurrency limit (`max_concurrent_checks`) is reached, additional due checks SHALL be queued. Queued checks SHALL be dispatched in priority order as running slots become available.

#### Scenario: Checks queued at concurrency limit
- **WHEN** `max_concurrent_checks` is 2 and 4 checks are due
- **THEN** 2 checks execute immediately and 2 are queued

#### Scenario: Queued check dispatched on completion
- **WHEN** a running check completes and a queued check exists
- **THEN** the queued check is dispatched from the front of the priority-sorted queue

### Requirement: No resource exhaustion
The executor and scheduler SHALL NOT spawn unbounded goroutines or containers regardless of the number of configured checks.

#### Scenario: Many checks configured
- **WHEN** 100 checks are configured with `max_concurrent_checks: 5`
- **THEN** at most 5 containers run concurrently
