## Why

The GitHub Actions release workflow is implemented at `.github/workflows/release.yml` with a legacy spec at `specs/github-actions-release.md`. Migrating to OpenSpec formalizes the CI/CD pipeline.

## What Changes

- Formalize the existing release workflow as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `github-actions-release`: CI/CD workflow that tests, cross-compiles, and publishes Foghorn binaries and Docker images for multiple architectures on push to main

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/github-actions-release.md` — legacy spec to be removed
- No code or workflow changes
