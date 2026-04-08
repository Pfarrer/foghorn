## Context

The Docker image build is handled by the `docker` job in `.github/workflows/release.yml`. It uses Docker Buildx with a matrix strategy for linux/amd64, linux/arm64, and linux/arm/v7. The `docker-manifest` job combines per-arch images into a multi-arch manifest. A `docker-compose.yml` example demonstrates how to run Foghorn as a container.

## Goals / Non-Goals

**Goals:**
- Document the Docker image build/publish and docker-compose example

**Non-Goals:**
- Changing build platforms, registry, or the docker-compose example

## Decisions

**Pre-built binary injection**: The Docker build receives the cross-compiled binary as a build arg (`BINARY`) and copies it into the image, rather than building Go code inside Docker.

**GHCR as registry**: Images are pushed to `ghcr.io/pfarrer/foghorn` using the automatic `GITHUB_TOKEN`.

**Docker Compose example**: Mounts Docker socket, config file as read-only volume, exposes status API on port 7676.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
