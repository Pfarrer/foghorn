## ADDED Requirements

### Requirement: Container directory structure
Each check container SHALL have its own subfolder under `containers/` containing at minimum a `Dockerfile` and a `README.md`.

#### Scenario: Container folder exists
- **WHEN** a container named `disk-check` is defined
- **THEN** `containers/disk-check/Dockerfile` and `containers/disk-check/README.md` exist

### Requirement: Automated build on change
The `container-release.yml` workflow SHALL trigger on pushes to `main` that modify files under `containers/**`. It SHALL detect which container folders changed and build only those.

#### Scenario: Changed container is built
- **WHEN** a file under `containers/openssl-check/` is modified and pushed to main
- **THEN** only the `openssl-check` container is built and pushed

#### Scenario: Unchanged containers skipped
- **WHEN** only `containers/openssl-check/` changed
- **THEN** `disk-check`, `http-check`, and other containers are not rebuilt

#### Scenario: No container changes
- **WHEN** a push to main modifies only non-container files
- **THEN** the build job is skipped entirely

### Requirement: Version tagging
Each container SHALL have a `VERSION` file with a semver string (`MAJOR.MINOR.PATCH`). The workflow SHALL tag images with the version and `latest`.

#### Scenario: Image tagged with version
- **WHEN** `containers/openssl-check/VERSION` contains `1.0.0`
- **THEN** the image is pushed as `ghcr.io/pfarrer/foghorn-openssl-check:1.0.0` and `ghcr.io/pfarrer/foghorn-openssl-check:latest`

#### Scenario: Invalid version rejected
- **WHEN** the `VERSION` file does not contain a valid semver string
- **THEN** the workflow fails with an error
