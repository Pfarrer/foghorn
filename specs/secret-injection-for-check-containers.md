# Secret Injection for Check Containers

## Category
security

## Description
Define a single secret storage and delivery model for check containers without storing cleartext secrets in config files.

## Usage Steps
1. Store secrets in an encrypted local secret store.
2. Reference secrets in check config by logical key (not cleartext value).
3. Run Foghorn and inject secrets into containers at runtime.

## Implementation Notes
- Use one secret provider model for all checks: local encrypted secret file plus a master key from OS keyring or env var at startup.
- Add a config reference format (for example `secret://smtp/password`) for check fields that require secrets.
- Resolve secret references only in memory at execution time.
- Inject secrets via ephemeral files or stdin mounted/read by the check container.
- Do not pass secrets as Docker CLI arguments.
- Keep secret values out of persisted check state and history files.
- Provide a secret management CLI to set, update, list keys, and delete keys without printing secret values.

## Acceptance Criteria
- [ ] Checks can reference secrets by key without embedding cleartext in YAML config.
- [ ] Secret values are encrypted at rest in the local secret store.
- [ ] Secret values are resolved only at runtime in memory.
- [ ] Secret values are delivered to check containers without Docker CLI secret arguments.
- [ ] Secret values are never written to state/history files.
- [ ] Secret management CLI supports create, rotate, and delete operations.
