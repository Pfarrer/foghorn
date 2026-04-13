## ADDED Requirements

### Requirement: Persistent memory config field
Each check SHALL support an optional `persistent_memory` boolean field. When `true`, the check container receives a persistent Docker volume mounted at `/run/foghorn/memory`.

#### Scenario: Persistent memory enabled
- **WHEN** a check has `persistent_memory: true`
- **THEN** the config is valid and the executor mounts a persistent volume

#### Scenario: Persistent memory disabled (default)
- **WHEN** a check omits `persistent_memory` or sets it to `false`
- **THEN** no persistent volume is mounted
