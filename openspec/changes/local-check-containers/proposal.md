## Why

The local check containers infrastructure is implemented with `containers/` directory and `.github/workflows/container-release.yml`, with a legacy spec at `specs/local-check-containers.md`. Migrating to OpenSpec formalizes the container build/release workflow.

## What Changes

- Formalize the local check containers feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `local-check-containers`: In-repo check container definitions with automated build and release triggered by changes to container folders

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/local-check-containers.md` — legacy spec to be removed
- No code or workflow changes
