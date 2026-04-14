## MODIFIED Requirements

### Requirement: JSON output format
Check containers SHALL output JSON to stdout with the following fields: `status` (string: "pass", "fail", "warn", "error", or "unknown"), `message` (string: human-readable description), `data` (optional object: structured metrics), `timestamp` (ISO 8601 timestamp), `duration_ms` (integer: execution duration).

The `error` status indicates the check could not finish due to a technical error in the check container (e.g., network failure, browser crash, missing dependency). It SHALL NOT be used to indicate a problem with the checked property — that is the purpose of `fail`.

#### Scenario: Valid JSON output parsed
- **WHEN** a container outputs valid JSON with all required fields
- **THEN** the executor extracts `status`, `message`, and `data` correctly

#### Scenario: Invalid JSON handled gracefully
- **WHEN** a container outputs non-JSON content
- **THEN** the executor records an error without crashing

#### Scenario: Missing JSON fields handled
- **WHEN** a container outputs JSON missing the `status` field
- **THEN** the executor records an error for that check

#### Scenario: Error status recorded
- **WHEN** a container outputs `{"status": "error", "message": "DNS resolution failed"}`
- **THEN** the executor records the result with status `error` and preserves the message

## ADDED Requirements

### Requirement: Error status semantics
The `error` status SHALL be treated as a distinct unhealthy state. It SHALL be persisted in the state log, counted separately in snapshot counts, and cause a non-zero exit in one-shot mode. The TUI SHALL render `error` with a distinct symbol and color from `fail` and `unknown`.

#### Scenario: Error status in snapshot counts
- **WHEN** a check's last status is `error`
- **THEN** the snapshot `counts.error` field is incremented for that check

#### Scenario: Error status in one-shot mode
- **WHEN** any check finishes with status `error` in one-shot mode
- **THEN** the daemon exits with a non-zero code

#### Scenario: Error status in TUI
- **WHEN** a check has last status `error`
- **THEN** the TUI displays a distinct symbol and color for that check row

#### Scenario: Error status persisted in state log
- **WHEN** a check completes with status `error`
- **THEN** a state log record is written with `status: "error"`
