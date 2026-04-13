## Why

Check containers are currently stateless — each run starts with a clean filesystem. This prevents checks from tracking changes across runs (e.g., detecting configuration drift, comparing metric baselines, or remembering previously seen alerts). Persistent memory allows check containers to store and retrieve data between executions.

## What Changes

- Add `persistent_memory` field to check config (boolean or volume config object)
- Mount a named Docker volume into check containers at a well-known path when enabled
- Volume persists across check runs, scoped per check name
- Check containers can read/write files in the persistent volume to carry state forward
- Volumes are created lazily on first run and cleaned up when checks are removed

## Capabilities

### New Capabilities
- `check-persistent-memory`: Docker volume-based persistent storage for check containers, scoped per check, with lifecycle management

### Modified Capabilities
- `define-config-yaml-format`: Add `persistent_memory` field to CheckConfig
- `docker-check-interface`: Add well-known mount path and environment variable for persistent memory volume

## Impact

- **config/types.go**: New `PersistentMemory` field on `CheckConfig`
- **executor/docker.go**: Volume creation and mounting logic in `Execute`
- **Config YAML**: New optional `persistent_memory` section per check
- **Check containers**: New env var `FOGHORN_PERSISTENT_DIR` pointing to mount path
- **Docker Engine**: Named volumes created per check (cleanup needed on check removal)
