## Context

The secret store (`secretstore/store.go`) encrypts secrets at rest using AES-GCM with a master key derived via Argon2. Secrets are referenced in check configs as `secret://<key>`. At runtime, the executor resolves these references and injects secret values via ephemeral files (`_FILE` env vars pointing to `/run/foghorn/secrets/`). A CLI (`foghorn-daemon secret`) manages the store with set, list, and delete commands.

## Goals / Non-Goals

**Goals:**
- Document the secret injection system as formal OpenSpec artifacts

**Non-Goals:**
- Adding remote secret providers, vault integration, or new reference formats

## Decisions

**AES-GCM encryption**: Secrets are encrypted at rest using AES-256-GCM with Argon2-derived keys. Each secret is stored as a JSON object with version, nonce, and ciphertext.

**Ephemeral file injection**: Secrets are written to temp files under `/run/foghorn/secrets/` and injected as `_FILE` env vars, avoiding Docker CLI arguments.

**No secrets in state logs**: Secret values are never written to state log files or history records.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
