## ADDED Requirements

### Requirement: Persistent memory volume injection
When a check has persistent memory enabled, the executor SHALL mount a named Docker volume at `/run/foghorn/memory` and inject the `FOGHORN_PERSISTENT_DIR` environment variable pointing to that path.

#### Scenario: Volume mount with persistent memory
- **WHEN** a check has `persistent_memory: true`
- **THEN** the container is created with a bind mount of volume `foghorn-memory-<check-name>` at `/run/foghorn/memory` and env var `FOGHORN_PERSISTENT_DIR=/run/foghorn/memory`

#### Scenario: No volume mount without persistent memory
- **WHEN** a check does not have persistent memory enabled
- **THEN** no additional volume mount or env var is added beyond the existing secret mount
