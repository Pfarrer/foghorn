## Context

The executor (`executor/docker.go`) injects environment variables into check containers and parses their JSON stdout output. The contract defines what containers receive (env vars) and must produce (JSON on stdout, exit codes).

## Goals / Non-Goals

**Goals:**
- Document the existing check interface contract as formal OpenSpec artifacts

**Non-Goals:**
- Changing env var names, output format, or error handling behavior

## Decisions

**Environment variable injection**: The executor injects `FOGHORN_CHECK_NAME`, `FOGHORN_CHECK_CONFIG`, `FOGHORN_ENDPOINT`, and `FOGHORN_TIMEOUT` as environment variables. Secrets are injected via `_FILE` suffixed vars pointing to ephemeral files.

**Stdout JSON parsing**: Container output is read from stdout and parsed as JSON. Non-zero exit codes are treated as failures regardless of JSON content.

**Graceful error handling**: Missing or invalid JSON output does not crash the executor; it results in an error being recorded for that check.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
