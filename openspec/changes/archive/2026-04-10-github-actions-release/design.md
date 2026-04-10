## Context

The release workflow (`.github/workflows/release.yml`) triggers on push to `main`. It runs tests first, then cross-compiles the daemon and TUI binaries for linux/amd64, linux/arm64, and linux/arm/v7 using Go's GOOS/GOARCH. Artifacts are uploaded, then Docker images are built per architecture and combined into a multi-arch manifest pushed to GHCR.

## Goals / Non-Goals

**Goals:**
- Document the existing release workflow as formal OpenSpec artifacts

**Non-Goals:**
- Adding new architectures, changing the trigger, or modifying build flags

## Decisions

**Go cross-compilation with CGO_ENABLED=0**: Builds are static binaries using `-ldflags="-s -w"` for size reduction.

**Matrix strategy**: Separate build jobs per architecture with explicit GOARM=7 for armhf.

**Docker multi-arch manifest**: Per-arch images are tagged individually, then combined via `docker buildx imagetools create` with `latest` and SHA tags.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
