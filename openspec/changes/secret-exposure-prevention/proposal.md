## Why

Secret exposure prevention is spec'd but not fully implemented. While the executor already uses ephemeral files (`_FILE` env vars) and the logger has redaction for debug output, the spec requires a centralized redaction package, process list protection, and endpoint serialization guards. Creating this as an OpenSpec change defines the requirements and implementation plan.

## What Changes

- Create a central `redact` package for consistent secret redaction across logger, executor, and endpoints
- Replace env var injection with file-based injection to avoid `ps` exposure
- Ensure endpoint serializers strip secret values
- Add integration tests for secret leak detection
- Add runtime guardrails that fail if secrets are about to be logged

## Capabilities

### New Capabilities
- `secret-exposure-prevention`: Centralized redaction ensuring secret values never appear in logs, process lists, or API responses

## Impact

- New `redact/` package with redaction functions
- Changes to `executor/docker.go` for file-based injection
- Changes to `internal/statusapi/` to strip secrets from responses
- Changes to `logger/` to use centralized redaction
- Tests for leak detection
- `specs/secret-exposure-prevention.md` — legacy spec to be removed
