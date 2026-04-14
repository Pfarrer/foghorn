## ADDED Requirements

### Requirement: Environment variable injection
The executor SHALL inject the following environment variables into every check container: `FOGHORN_CHECK_NAME` (check name), `FOGHORN_CHECK_CONFIG` (JSON string of check metadata, when present), `FOGHORN_ENDPOINT` (target endpoint, when configured), and `FOGHORN_TIMEOUT` (parsed timeout duration).

#### Scenario: All env vars injected
- **WHEN** a check is executed with name "api-health", endpoint "https://example.com", and timeout "30s"
- **THEN** the container receives `FOGHORN_CHECK_NAME=api-health`, `FOGHORN_ENDPOINT=https://example.com`, and `FOGHORN_TIMEOUT=30s`

#### Scenario: Optional env vars omitted when unset
- **WHEN** a check has no endpoint or metadata configured
- **THEN** `FOGHORN_ENDPOINT` and `FOGHORN_CHECK_CONFIG` are not injected

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

### Requirement: Exit code semantics
Non-zero exit codes SHALL be treated as check failures. Zero exit code indicates success. The JSON `status` field takes precedence when present and exit code is zero.

#### Scenario: Non-zero exit code
- **WHEN** a container exits with code 1
- **THEN** the check is treated as failed

#### Scenario: Zero exit code with pass status
- **WHEN** a container exits with code 0 and outputs `{"status": "pass"}`
- **THEN** the check is recorded as passed

### Requirement: Secret injection via files
Secrets SHALL be injected as environment variables suffixed with `_FILE` pointing to ephemeral files at `/run/foghorn/secrets/`. The secret value is written to the file, not passed as a direct env var.

#### Scenario: Secret file injected
- **WHEN** a check has an env var with a secret reference
- **THEN** a corresponding `_FILE` env var points to a file containing the secret value
