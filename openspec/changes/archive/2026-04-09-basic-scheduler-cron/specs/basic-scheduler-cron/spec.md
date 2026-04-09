## ADDED Requirements

### Requirement: Cron schedule parsing
The scheduler SHALL parse standard 5-field cron expressions (minute, hour, day of month, month, day of week) from check configurations. Invalid cron expressions SHALL produce an error when adding the check.

#### Scenario: Valid cron expression
- **WHEN** a check is added with schedule `*/5 * * * *`
- **THEN** the check is added successfully with a calculated next run time

#### Scenario: Invalid cron expression
- **WHEN** a check is added with schedule `invalid-cron`
- **THEN** `AddCheck` returns an error indicating the cron expression failed to parse

### Requirement: Next run time calculation
The scheduler SHALL calculate the next execution time for each cron-scheduled check. The next run time SHALL be the earliest future time matching the cron expression in the configured time zone.

#### Scenario: Next run is in the future
- **WHEN** a cron check with schedule `0 0 * * *` (midnight daily) is added at 10:00
- **THEN** next run is set to midnight tonight (or tomorrow if past midnight)

### Requirement: Cron check execution
The scheduler SHALL trigger execution of cron-scheduled checks when the current time reaches or passes their `NextRun`. After execution, the `NextRun` SHALL be updated to the following cron match.

#### Scenario: Check fires on schedule
- **WHEN** the ticker fires and current time equals the check's `NextRun`
- **THEN** the check is executed

#### Scenario: Next run updated after execution
- **WHEN** a cron check completes execution
- **THEN** its `NextRun` is recalculated to the next cron match

### Requirement: Time zone support
The scheduler SHALL operate in the configured time location. All time comparisons and cron evaluations SHALL use the configured zone.

#### Scenario: Time zone respected
- **WHEN** the scheduler is configured with `America/New_York`
- **THEN** cron expressions are evaluated against the New York time zone

### Requirement: Scheduler lifecycle
The scheduler SHALL support `Start` and `Stop` operations. `Start` begins the ticker loop; `Stop` halts it cleanly.

#### Scenario: Start begins ticking
- **WHEN** `Start` is called
- **THEN** the scheduler begins evaluating due checks on each tick

#### Scenario: Stop halts execution
- **WHEN** `Stop` is called
- **THEN** the ticker is stopped and the run loop exits
