# Docker Image Build and Publish

## Category
integration

## Description
Extend the GitHub Actions workflow to build and publish a multi-architecture Docker image to GitHub Container Registry (ghcr.io)

## Usage Steps
1. Push code to main branch
2. GitHub Actions triggers the workflow
3. Workflow builds multi-architecture Docker image (amd64, arm64, arm/v7)
4. Image is pushed to GitHub Container Registry
5. Image is tagged according to release versioning policy
6. Pull and run image using docker compose example

## Implementation Notes
- Extend `.github/workflows/release.yml` workflow file
- Add new job 'docker' that depends on test job
- Use Docker Buildx for multi-platform builds
- Configure platforms: linux/amd64, linux/arm64, linux/arm/v7
- Authenticate with GitHub Container Registry using GITHUB_TOKEN
- Push images to ghcr.io/pfarrer/foghorn
- Tag images according to release versioning policy
- Use docker/metadata-action to generate tags and labels
- Set proper Dockerfile context and build args
- Ensure image uses non-root user if applicable
- Use layer caching to speed up builds
- Create example docker-compose.yml file demonstrating basic usage
- Example should include volume mounts for config and any necessary environment variables

## Acceptance Criteria
- [x] Docker job runs after tests pass
- [x] Image builds successfully for all three platforms
- [x] Image is pushed to GitHub Container Registry
- [x] Image tags follow the release versioning policy
- [x] Build logs show multi-platform build progress
- [x] Workflow fails if Docker build or push fails
- [x] Example docker-compose.yml is provided demonstrating general usage
