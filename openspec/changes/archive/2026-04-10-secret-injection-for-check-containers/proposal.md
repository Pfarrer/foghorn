## Why

Secret injection is implemented in `secretstore/store.go` and `executor/docker.go` with a legacy spec at `specs/secret-injection-for-check-containers.md`. Migrating to OpenSpec formalizes the secret storage and delivery model.

## What Changes

- Formalize the secret injection feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `secret-injection-for-check-containers`: Encrypted local secret store with `secret://` config references, runtime-only resolution, and ephemeral file injection into check containers

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/secret-injection-for-check-containers.md` — legacy spec to be removed
- No code changes
