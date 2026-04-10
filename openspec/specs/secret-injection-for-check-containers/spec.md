# secret-injection-for-check-containers Specification

## Purpose
TBD - created by archiving change secret-injection-for-check-containers. Update Purpose after archive.
## Requirements
### Requirement: Secret reference format
Check configs SHALL reference secrets using the `secret://<key>` format. The executor SHALL detect these references at runtime and resolve them to their actual values.

#### Scenario: Secret reference parsed
- **WHEN** a config value is `secret://smtp/password`
- **THEN** the reference key is `smtp/password`

#### Scenario: Non-secret values unchanged
- **WHEN** a config value is `plain-text-value`
- **THEN** it is passed through without secret resolution

### Requirement: Encrypted at-rest storage
Secrets SHALL be stored encrypted in a local file using AES-256-GCM. The master key SHALL be derived via Argon2 from a user-provided passphrase.

#### Scenario: Secret encrypted on write
- **WHEN** a secret is stored via the CLI
- **THEN** it is written as encrypted ciphertext to the store file

#### Scenario: Secret decrypted on read
- **WHEN** a `secret://` reference is resolved at runtime
- **THEN** the ciphertext is decrypted and the plaintext value is available in memory only

### Requirement: Runtime-only resolution
Secret values SHALL be resolved only in memory at execution time. They SHALL NOT be persisted to state logs, history files, or any other persistent storage.

#### Scenario: Secrets not in state logs
- **WHEN** a check using secrets completes and results are recorded
- **THEN** the state log contains no secret values

#### Scenario: Secrets not in history
- **WHEN** check execution history is displayed
- **THEN** secret values are not included

### Requirement: Ephemeral file injection
Secrets SHALL be injected into check containers via `_FILE` environment variables pointing to ephemeral files. Secrets SHALL NOT be passed as Docker CLI arguments.

#### Scenario: Secret file env var
- **WHEN** a check env var `SMTP_PASSWORD` references `secret://smtp/password`
- **THEN** `SMTP_PASSWORD_FILE` is set to an ephemeral file path containing the secret

#### Scenario: No Docker CLI secret args
- **WHEN** secrets are injected
- **THEN** they are never passed as Docker command-line arguments

### Requirement: Secret management CLI
The daemon SHALL provide a `secret` subcommand supporting create/set, list keys, and delete operations. Secret values SHALL NOT be printed by list operations.

#### Scenario: Set a secret
- **WHEN** `foghorn-daemon secret set smtp/password` is run
- **THEN** the user is prompted for the value and it is stored encrypted

#### Scenario: List keys
- **WHEN** `foghorn-daemon secret list` is run
- **THEN** only key names are displayed, not values

#### Scenario: Delete a secret
- **WHEN** `foghorn-daemon secret delete smtp/password` is run
- **THEN** the key is removed from the store

