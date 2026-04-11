## ADDED Requirements

### Requirement: CLI flag
The daemon SHALL accept a `--one-shot` flag to enable one-shot mode. The flag SHALL be combinable with other flags like `--config`, `--log-level`, and `--dry-run`.

#### Scenario: One-shot flag accepted
- **WHEN** `--one-shot` is provided
- **THEN** one-shot mode is enabled

#### Scenario: One-shot with config
- **WHEN** `--one-shot -c config.yaml` is provided
- **THEN** the config is loaded and one-shot mode runs

### Requirement: Single execution per check
In one-shot mode, each enabled check SHALL be executed exactly once. Recurring scheduling is bypassed.

#### Scenario: Each check runs once
- **WHEN** 5 checks are configured and enabled
- **THEN** each check runs exactly one time

#### Scenario: Disabled checks skipped
- **WHEN** a check has `enabled: false`
- **THEN** it is not executed in one-shot mode

### Requirement: Process exit after completion
After all checks have been executed, the daemon SHALL exit. The scheduler tick loop is not started.

#### Scenario: Process exits after runs
- **WHEN** all checks have completed
- **THEN** the process exits

### Requirement: Existing evaluator logic
One-shot mode SHALL use the same result evaluation logic as normal scheduler runs (pass/fail/warn determination).

#### Scenario: Evaluator used
- **WHEN** a check completes in one-shot mode
- **THEN** the result is evaluated using the existing evaluator

### Requirement: Exit code reflecting results
The process SHALL exit with code 0 when all checks pass or warn. Exit with code 1 when any check fails or results in unknown status.

#### Scenario: All pass/warn
- **WHEN** all check results are `pass` or `warn`
- **THEN** the process exits with code 0

#### Scenario: Any fail
- **WHEN** any check result is `fail` or `unknown`
- **THEN** the process exits with code 1
