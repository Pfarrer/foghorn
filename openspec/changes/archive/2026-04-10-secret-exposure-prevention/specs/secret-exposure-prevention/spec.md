## ADDED Requirements

### Requirement: Central redaction package
A new `redact` package SHALL provide a `Sanitize(input string, secrets []string) string` function that replaces all known secret values and common patterns with `[REDACTED]`. This package SHALL be used by the logger, executor, and status API.

#### Scenario: Secrets redacted
- **WHEN** input contains a known secret value "super-secret-password"
- **THEN** it is replaced with `[REDACTED]`

#### Scenario: Auth patterns redacted
- **WHEN** input contains `Authorization: Bearer token123`
- **THEN** it becomes `Authorization: Bearer [REDACTED]`

#### Scenario: Common patterns redacted
- **WHEN** input contains `password=secret123` or `token=abc`
- **THEN** the values are redacted

### Requirement: No secrets in logs
Resolved secret values SHALL NOT appear in application logs at any log level.

#### Scenario: Log contains no secrets
- **WHEN** a check runs with secret references
- **THEN** no log lines contain any resolved secret values

#### Scenario: Secret leak test fails
- **WHEN** a test detects a secret in log output
- **THEN** the test fails

### Requirement: No secrets in process lists
Secrets SHALL NOT appear in Linux process argument lists. Ephemeral file paths (not secret values) are acceptable in environment variables shown by `ps`.

#### Scenario: No secret values in ps
- **WHEN** `ps aux` is run while Foghorn is running
- **THEN** no resolved secret values appear in the command line or environment

#### Scenario: File paths acceptable
- **WHEN** an env var is `SMTP_PASSWORD_FILE=/run/foghorn/secrets/xyz`
- **THEN** it is acceptable in `ps` output (contains path, not secret)

### Requirement: No secrets in endpoint responses
The status API endpoint responses SHALL NOT include secret values from check configurations.

#### Scenario: API response sanitized
- **WHEN** a client fetches status from the API
- **THEN** the response contains no secret values from check env vars

### Requirement: Consistent redaction across all paths
Redaction SHALL be applied consistently: logger, error paths, debug output, and status API.

#### Scenario: All paths use redaction
- **WHEN** any output (log, API, debug) includes potentially secret data
- **THEN** it is sanitized by the central redact package

### Requirement: Integration tests for leak detection
Automated tests SHALL verify that no secret leaks occur in logs or endpoint responses.

#### Scenario: Leak detection test
- **WHEN** an integration test runs checks and inspects logs and API output
- **THEN** it fails if any secret value is detected
