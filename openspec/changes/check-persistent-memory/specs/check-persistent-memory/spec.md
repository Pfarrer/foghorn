## ADDED Requirements

### Requirement: Persistent volume mounting
When a check has `persistent_memory: true`, the executor SHALL create (if needed) a named Docker volume `foghorn-memory-<check-name>` and mount it at `/run/foghorn/memory` inside the container with read-write access.

#### Scenario: First execution with persistent memory
- **WHEN** a check named "ssl-expiry" has `persistent_memory: true` and no volume `foghorn-memory-ssl-expiry` exists
- **THEN** the executor creates the volume before starting the container and mounts it at `/run/foghorn/memory`

#### Scenario: Subsequent execution reuses volume
- **WHEN** a check named "ssl-expiry" has `persistent_memory: true` and volume `foghorn-memory-ssl-expiry` already exists
- **THEN** the executor mounts the existing volume without recreating it

#### Scenario: Persistent memory disabled
- **WHEN** a check has no `persistent_memory` field or it is set to `false`
- **THEN** no volume is created or mounted

### Requirement: Persistent directory environment variable
When a check has `persistent_memory: true`, the executor SHALL inject the environment variable `FOGHORN_PERSISTENT_DIR=/run/foghorn/memory` into the container.

#### Scenario: Env var injected with persistent memory
- **WHEN** a check has `persistent_memory: true`
- **THEN** the container receives `FOGHORN_PERSISTENT_DIR=/run/foghorn/memory`

#### Scenario: Env var omitted without persistent memory
- **WHEN** a check has `persistent_memory: false` or unset
- **THEN** `FOGHORN_PERSISTENT_DIR` is not injected

### Requirement: Volume naming convention
Persistent volumes SHALL be named with the prefix `foghorn-memory-` followed by the sanitized check name (lowercase, non-alphanumeric replaced with `-`).

#### Scenario: Standard check name
- **WHEN** a check is named "disk-usage"
- **THEN** the volume is named `foghorn-memory-disk-usage`

#### Scenario: Check name with special characters
- **WHEN** a check is named "SSL_Check (prod)"
- **THEN** the volume is named `foghorn-memory-ssl-check--prod-`

### Requirement: Volume isolation
Each check SHALL have its own volume. Volumes SHALL NOT be shared between different checks.

#### Scenario: Two checks with persistent memory
- **WHEN** checks "check-a" and "check-b" both have `persistent_memory: true`
- **THEN** two separate volumes `foghorn-memory-check-a` and `foghorn-memory-check-b` are created and mounted independently
