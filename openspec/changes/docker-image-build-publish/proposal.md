## Why

The Docker image build and publish workflow is part of `.github/workflows/release.yml` (the `docker` and `docker-manifest` jobs) with a legacy spec at `specs/docker-image-build-publish.md`. An example `docker-compose.yml` is also provided. Migrating to OpenSpec formalizes this.

## What Changes

- Formalize the Docker image build/publish workflow and docker-compose example as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `docker-image-build-publish`: Multi-architecture Docker image build and publish to GHCR via GitHub Actions, with a docker-compose example for deployment

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/docker-image-build-publish.md` — legacy spec to be removed
- No code or workflow changes
