# GitHub Actions Release Workflow

## Category
integration

## Description
GitHub Actions workflow that automatically builds and releases Foghorn binaries for multiple architectures on every push to main branch

## Usage Steps
1. Push code to main branch
2. GitHub Actions triggers automatically
3. Workflow builds Foghorn binaries for arm64, armhf, and x64
4. Release artifacts are created and attached to the workflow run

## Implementation Notes
- Create `.github/workflows/release.yml` workflow file
- Configure workflow to trigger on push to main branch
- Use Go cross-compilation to build for target architectures:
  - linux/amd64 (x64)
  - linux/arm64
  - linux/arm/v7 (armhf)
- Set appropriate GOARCH and GOARM environment variables during build
- Create release artifacts for each architecture
- Use GitHub Actions artifact upload functionality
- Ensure builds are reproducible (disable module stripping or use consistent timestamps)
- Consider using GoReleaser if additional functionality needed

## Acceptance Criteria
- [x] Workflow triggers on every push to main branch
- [x] All tests must pass before building artifacts
- [x] Builds successfully complete for linux/amd64
- [x] Builds successfully complete for linux/arm64
- [x] Builds successfully complete for linux/arm/v7
- [x] Release artifacts are uploaded for each architecture
- [x] Artifact naming convention includes architecture (e.g., foghorn-linux-amd64, foghorn-linux-arm64, foghorn-linux-armhf)
- [x] Workflow fails if any architecture build fails
- [x] Build logs are clear and informative
