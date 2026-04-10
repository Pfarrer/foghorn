## ADDED Requirements

### Requirement: Concurrency limit configuration
The system SHALL support a `max_concurrent_checks` config property (int, optional). When set to a positive value, at most that many check containers SHALL run simultaneously. When unset or 0, concurrency SHALL be unlimited. Negative values SHALL be rejected at config load time.

#### Scenario: Concurrency limit enforced
- **WHEN** `max_concurrent_checks` is set to 3 and 5 checks are due simultaneously
- **THEN** at most 3 checks execute concurrently; the remaining 2 are queued

#### Scenario: Unlimited concurrency
- **WHEN** `max_concurrent_checks` is unset or 0
- **THEN** all due checks execute immediately without concurrency limiting

#### Scenario: Negative value rejected
- **WHEN** `max_concurrent_checks` is set to a negative value
- **THEN** config loading returns an error

### Requirement: Priority ordering of due checks
When multiple checks are due at the same time, the system SHALL prioritize checks with longer scheduling intervals (next planned run farther in the future). For equal intervals, alphabetical check name SHALL be used as a stable tie-breaker.

#### Scenario: Longer interval runs first
- **WHEN** check A has interval 5m and check B has interval 1h, and both are due
- **THEN** check B (longer interval) executes first

#### Scenario: Equal interval tie-breaker
- **WHEN** check "alpha" and check "beta" both have interval 5m and are due
- **THEN** "alpha" executes before "beta" based on alphabetical name ordering

#### Scenario: Single due check
- **WHEN** only one check is due
- **THEN** it executes immediately regardless of priority logic

### Requirement: Immediate run for interval checks on empty state
When the scheduler has no prior execution state, interval-based checks SHALL run immediately without an initial delay.

#### Scenario: First run with interval check
- **WHEN** the scheduler starts with empty state and a check has interval 10m
- **THEN** the check runs immediately

#### Scenario: Subsequent runs follow interval
- **WHEN** the interval check has already run once
- **THEN** the next run is scheduled 10m after the previous run

### Requirement: Cron checks wait for scheduled time
Cron-scheduled checks SHALL NOT run immediately, even on empty state. They SHALL wait for their next scheduled time.

#### Scenario: First run with cron check
- **WHEN** the scheduler starts with empty state and a check has cron schedule `*/5 * * * *`
- **THEN** the check does not run immediately; it waits for the next cron match

#### Scenario: Cron check fires on schedule
- **WHEN** the current time matches the cron expression
- **THEN** the cron check executes
