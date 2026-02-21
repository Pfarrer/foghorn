# Secret Exposure Prevention in Logs, Process Lists, and Endpoints

## Category
security

## Description
Guarantee that secret values are never exposed in logs, Linux process lists, or externally exposed service endpoints.

## Usage Steps
1. Configure checks with secret references.
2. Run Foghorn in normal daemon mode.
3. Review logs, process metadata, and external endpoint outputs for redaction behavior.

## Implementation Notes
- Add a central redaction package used by logger, executor, and endpoint serializers.
- Enforce a default-deny output model: secrets are never included unless explicitly marked safe.
- Redact known secret patterns and all resolved secret values before any log write.
- Avoid placing secrets in command arguments or environment when they can appear in `ps` output.
- Ensure endpoint payloads include only non-sensitive check metadata and status.
- Add integration tests that inspect logs and API responses for leaked secret values.
- Add runtime guardrails that fail execution if a secret is about to be logged or serialized.

## Acceptance Criteria
- [ ] No resolved secret value appears in application logs.
- [ ] No resolved secret value appears in Linux process argument lists.
- [ ] No resolved secret value appears in exported endpoint responses.
- [ ] Redaction is applied consistently across all log levels and error paths.
- [ ] Automated tests fail when a secret leak is detected in logs or endpoint output.
