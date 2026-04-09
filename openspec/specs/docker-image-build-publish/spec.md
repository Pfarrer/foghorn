## ADDED Requirements

### Requirement: Multi-architecture Docker build
The release workflow SHALL build Docker images for linux/amd64, linux/arm64, and linux/arm/v7 using Docker Buildx. Builds SHALL run after tests pass and after the corresponding binary artifact is downloaded.

#### Scenario: All platforms build
- **WHEN** the docker build job runs
- **THEN** images are built and pushed for all three platforms to ghcr.io

#### Scenario: Build failure stops workflow
- **WHEN** a Docker build fails for any platform
- **THEN** the workflow fails and the manifest job does not proceed

### Requirement: Docker image tagging
Per-architecture images SHALL be tagged as `ghcr.io/pfarrer/foghorn:linux-<arch>-<sha>`. The multi-arch manifest SHALL be tagged as `latest` and `<sha>`.

#### Scenario: Image tags
- **WHEN** a linux/arm64 image is built on commit abc123
- **THEN** it is pushed as `ghcr.io/pfarrer/foghorn:linux-arm64-abc123`

#### Scenario: Manifest tags
- **WHEN** the manifest is created
- **THEN** `ghcr.io/pfarrer/foghorn:latest` and `ghcr.io/pfarrer/foghorn:abc123` resolve to the correct architecture

### Requirement: Docker Compose example
A `docker-compose.yml` SHALL be provided demonstrating how to run Foghorn as a container with Docker socket access, config volume mount, and status API exposure.

#### Scenario: Example is valid
- **WHEN** `docker-compose up` is run with the example file
- **THEN** Foghorn starts with the provided config and Docker socket access

#### Scenario: Status API exposed
- **WHEN** the example is running
- **THEN** the status API is accessible on port 7676
