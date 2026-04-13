## Context

Foghorn runs check containers as ephemeral Docker workloads. Each execution creates a fresh container with no filesystem persistence. This is sufficient for stateless checks (HTTP ping, disk usage) but prevents checks that need to compare current state to previous runs (e.g., detecting config drift, tracking metric trends, deduplicating alerts).

The executor already uses Docker Engine API with host binds for secrets. Extending it to mount named Docker volumes is a natural fit.

## Goals / Non-Goals

**Goals:**
- Allow check containers to persist files across executions
- Scope persistent storage per check name (isolated volumes)
- Simple opt-in config: `persistent_memory: true`
- Provide a well-known mount path and env var inside containers
- Lazy volume creation on first execution

**Non-Goals:**
- Volume size limits or quotas (Docker native limits suffice for now)
- Cross-check data sharing
- Volume encryption beyond what Docker provides
- Automatic volume cleanup on check removal (separate change)
- Backup or migration of persistent memory contents

## Decisions

### 1. Named Docker volumes per check

**Decision**: Use named Docker volumes with naming convention `foghorn-memory-<check-name>`.

**Rationale**: Named volumes survive container removal (the executor already removes containers after each run). They are managed by Docker, require no host path management, and are portable across Docker environments.

**Alternatives considered**:
- **Host bind mounts**: Require managing host directories, permissions, and cleanup. Less portable.
- **Docker tmpfs**: Not persistent across container runs.

### 2. Mount path `/run/foghorn/memory`

**Decision**: Mount persistent volumes at `/run/foghorn/memory` inside containers.

**Rationale**: Follows the existing `/run/foghorn/secrets` convention. Predictable path for check container authors. Communicated via `FOGHORN_PERSISTENT_DIR` env var.

### 3. Config format: `persistent_memory` boolean

**Decision**: `persistent_memory: true` on the check config. Future expansion to an object with options (size limit, path override) is possible without breaking changes via YAML unmarshaling.

**Rationale**: Minimal config surface. Most checks won't need this, so default-off keeps config clean.

### 4. Volume lifecycle: lazy creation

**Decision**: Create volumes on first check execution if they don't exist. Do not delete volumes when checks are removed from config.

**Rationale**: Simple to implement. Data loss from accidental config removal is avoided. Explicit cleanup can be added later.

## Risks / Trade-offs

- **Unbounded disk growth** → Mitigation: Docker volumes are inspectable; admins can monitor. Future: add size limits or TTL-based cleanup.
- **Volume name collisions with manual Docker use** → Mitigation: `foghorn-memory-` prefix makes ownership clear.
- **Check name changes create orphaned volumes** → Mitigation: Document that renaming a check creates a new volume. Old volume retained safely.
- **Concurrent check runs writing to same volume** → Mitigation: Same check name shouldn't run concurrently (scheduler enforces this already).
