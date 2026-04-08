## Why

The verify image availability flag is implemented in `internal/daemon/app.go` with tests in `app_test.go`, and a legacy spec at `specs/verify-image-availability-flag.md`. Migrating to OpenSpec formalizes this feature.

## What Changes

- Formalize the existing image verification feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `verify-image-availability-flag`: CLI flag (`-i`/`--verify-image-availability`) that validates all Docker images in the config are available locally before starting the scheduler

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/verify-image-availability-flag.md` — legacy spec to be removed
- No code changes
