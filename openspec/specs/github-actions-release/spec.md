# github-actions-release Specification

## Purpose
TBD - created by archiving change github-actions-release. Update Purpose after archive.
## Requirements
### Requirement: Workflow trigger
The release workflow SHALL trigger on every push to the `main` branch.

#### Scenario: Push to main triggers workflow
- **WHEN** code is pushed to the `main` branch
- **THEN** the release workflow is triggered automatically

### Requirement: Pre-build test gate
All tests SHALL pass before any build artifacts are created. If tests fail, the workflow SHALL stop without building.

#### Scenario: Tests pass
- **WHEN** `go test ./...` succeeds
- **THEN** the build jobs proceed

#### Scenario: Tests fail
- **WHEN** `go test ./...` fails
- **THEN** the workflow stops and no artifacts are produced

### Requirement: Cross-architecture builds
The workflow SHALL build both the daemon and TUI binaries for linux/amd64, linux/arm64, and linux/arm/v7 (armhf). Builds SHALL use `CGO_ENABLED=0` and `-ldflags="-s -w"`.

#### Scenario: All architectures build
- **WHEN** the build job runs
- **THEN** binaries are produced for all three architectures

#### Scenario: Build failure stops workflow
- **WHEN** any architecture build fails
- **THEN** the workflow fails and Docker image jobs do not proceed

### Requirement: Artifact upload
Built binaries SHALL be uploaded as GitHub Actions artifacts with naming convention `foghorn-daemon-linux-<arch>` and `foghorn-tui-linux-<arch>`.

#### Scenario: Artifacts uploaded
- **WHEN** a build succeeds for linux/amd64
- **THEN** `foghorn-daemon-linux-amd64` and `foghorn-tui-linux-amd64` artifacts are available

### Requirement: Docker image build and push
Per-architecture Docker images SHALL be built and pushed to GitHub Container Registry (ghcr.io) tagged with OS, architecture, and commit SHA.

#### Scenario: Docker image pushed
- **WHEN** a binary build succeeds for linux/arm64
- **THEN** `ghcr.io/pfarrer/foghorn:linux-arm64-<sha>` is pushed

### Requirement: Multi-arch manifest
A multi-arch Docker manifest SHALL be created combining all architecture-specific images, tagged as `latest` and with the commit SHA.

#### Scenario: Manifest created
- **WHEN** all per-arch Docker images are pushed
- **THEN** `ghcr.io/pfarrer/foghorn:latest` and `ghcr.io/pfarrer/foghorn:<sha>` resolve to the correct architecture on `docker pull`

