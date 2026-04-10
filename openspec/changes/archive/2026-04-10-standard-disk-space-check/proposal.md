## Why

The disk-check container exists as a fully implemented Docker container in `containers/disk-check/` with a legacy spec in `specs/standard-disk-space-check.md`. Migrating to OpenSpec provides structured artifacts for ongoing maintenance and consistency with the project's change management process.

## What Changes

- Formalize the existing disk-check container as an OpenSpec change with full artifacts
- Port the legacy spec into the OpenSpec spec format under `openspec/specs/disk-check/`
- No functional changes to the container or its behavior

## Capabilities

### New Capabilities
- `disk-check`: Filesystem disk space and inode monitoring via Docker container — validates usage against percentage or byte-based thresholds with optional inode checking

### Modified Capabilities
<!-- None — this is a migration of an existing implemented feature -->

## Impact

- `openspec/specs/` — new spec directory and spec.md for disk-check
- `specs/standard-disk-space-check.md` — legacy spec to be removed
- No code changes to `containers/disk-check/` or any Go packages
