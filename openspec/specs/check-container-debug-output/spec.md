# check-container-debug-output Specification

## Purpose
Define how check container debug output is enabled, redacted, and truncated.
## Requirements
### Requirement: Debug output modes
The system SHALL support three debug output modes: `off` (no output), `on_failure` (log on check failure or timeout), and `always` (log every check run). The mode SHALL be configurable globally via `check_container_debug_output` and overridable per check. Default SHALL be `off`.

#### Scenario: Off mode
- **WHEN** mode is `off`
- **THEN** no container output is logged regardless of check result

#### Scenario: On-failure mode with passing check
- **WHEN** mode is `on_failure` and the check passes
- **THEN** no container output is logged

#### Scenario: On-failure mode with failing check
- **WHEN** mode is `on_failure` and the check fails
- **THEN** container output is logged at debug level

#### Scenario: On-failure mode with timeout
- **WHEN** mode is `on_failure` and the check times out
- **THEN** container output is logged at debug level

#### Scenario: Always mode
- **WHEN** mode is `always` and the check passes
- **THEN** container output is logged at debug level

#### Scenario: Per-check override
- **WHEN** global mode is `off` but a check has `check_container_debug_output: "always"`
- **THEN** container output is logged for that check

### Requirement: Secret redaction
Container output SHALL be redacted before logging. Known secret values SHALL be replaced with `[REDACTED]`. Common auth patterns (Authorization headers, passwords, Bearer tokens) SHALL also be redacted.

#### Scenario: Secret values redacted
- **WHEN** container output contains a known secret value
- **THEN** it is replaced with `[REDACTED]` in the logged output

#### Scenario: Auth patterns redacted
- **WHEN** container output contains `Authorization: Bearer <token>` or `password=secret`
- **THEN** the sensitive parts are replaced with `[REDACTED]`

#### Scenario: No secret leaks
- **WHEN** debug output is logged
- **THEN** no secret values appear in the log output

### Requirement: Output truncation
Logged container output SHALL be truncated to a configurable maximum character limit (default: 4096).

#### Scenario: Output truncated
- **WHEN** container output exceeds the max character limit
- **THEN** the logged output is truncated at the limit

#### Scenario: Short output not truncated
- **WHEN** container output is within the max character limit
- **THEN** the full output is logged
