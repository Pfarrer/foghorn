## Context

Currently, secret values are handled inconsistently: the executor uses `_FILE` env vars (good), debug output is redacted (partial), but the logger has no centralized redaction for main logs, and the status API may expose secrets if check configs include them. The spec requires a unified approach.

## Goals / Non-Goals

**Goals:**
- Create a central `redact` package with functions to sanitize strings
- Ensure secrets never appear in application logs
- Ensure secrets never appear in process arguments (`ps aux`)
- Ensure endpoint responses never include secret values
- Add integration tests that verify no secret leaks

**Non-Goals:**
- Preventing secrets from being in container environment (that's handled by `_FILE`)
- Preventing secrets in Foghorn's own config file (user responsibility)

## Decisions

**Central redact package**: A new `redact/` package with `Sanitize(input string, secrets []string) string` that replaces all known secret values and common patterns (auth headers, passwords, tokens).

**File-based injection already in place**: The executor already uses `*_FILE` env vars, which is correct. No change needed there.

**Status API sanitization**: The status API already returns check metadata, not secrets. If check configs might include secret values in metadata fields, those fields must be stripped.

**Runtime guardrails**: Add a check before logging any value known to contain secrets (e.g., log lines that include check env vars). Fail the operation or replace with `[SECRET]` if detected.

## Risks / Trade-offs

- [Performance] → Redaction on every log write adds latency; mitigated by only redacting known secret values, not scanning all output
- [Missed patterns] → New secret patterns may appear; mitigated by parameterized logging and runtime guardrails
