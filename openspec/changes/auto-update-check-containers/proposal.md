## Why

Auto-update for check containers is spec'd but not yet implemented (all acceptance criteria unchecked). Creating this as an OpenSpec change defines the requirements and implementation plan.

## What Changes

- Define the auto-update feature with full OpenSpec artifacts
- Add config fields for enabling and scheduling auto-update
- Implement periodic Docker image pulls for configured check containers
- Log update attempts without blocking normal check execution

## Capabilities

### New Capabilities
- `auto-update-check-containers`: Periodic automatic pulling of latest check container images on a configurable schedule, with graceful failure handling

### Modified Capabilities
<!-- None -->

## Impact

- `config/types.go` — new config fields (`auto_update_containers`, `auto_update_schedule`)
- `config/loader.go` — validation for new fields
- `internal/daemon/app.go` — scheduling auto-update job
- `executor/docker.go` or new package — image pull logic
- `specs/auto-update-check-containers.md` — legacy spec to be removed
