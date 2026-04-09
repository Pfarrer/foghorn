## ADDED Requirements

### Requirement: Interval parsing
The scheduler SHALL parse interval strings in the format `<number><unit>` where unit is `s` (seconds), `m` (minutes), `h` (hours), or `d` (days). The number MUST be positive. Invalid intervals SHALL produce an error.

#### Scenario: Valid interval parsed
- **WHEN** `parseInterval("5m")` is called
- **THEN** it returns a 5-minute duration

#### Scenario: All supported units
- **WHEN** intervals `30s`, `2h`, `1d` are parsed
- **THEN** they return 30 seconds, 2 hours, and 24 hours respectively

#### Scenario: Invalid interval
- **WHEN** `parseInterval("abc")` or `parseInterval("-5m")` or `parseInterval("5x")` is called
- **THEN** an error is returned with a descriptive message

#### Scenario: Empty interval
- **WHEN** `parseInterval("")` is called
- **THEN** an error is returned indicating interval cannot be empty

### Requirement: Interval check execution
Interval-scheduled checks SHALL be triggered when the current time reaches or passes their `NextRun`. After execution, `NextRun` SHALL be set to the previous run time plus the interval duration.

#### Scenario: Check fires on interval
- **WHEN** a check has interval `5m` and the current time reaches `NextRun`
- **THEN** the check executes and `NextRun` is advanced by 5 minutes

#### Scenario: Immediate first run on empty state
- **WHEN** an interval check is added with no prior state
- **THEN** `NextRun` is set to the current time (runs immediately)

### Requirement: Mixed cron and interval schedules
The scheduler SHALL support both cron and interval checks in the same configuration without conflict.

#### Scenario: Mixed schedule types
- **WHEN** one check uses `cron: "*/5 * * * *"` and another uses `interval: "1m"`
- **THEN** both checks are managed by the scheduler and fire according to their respective schedules
