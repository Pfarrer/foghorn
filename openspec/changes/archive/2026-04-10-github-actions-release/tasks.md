## 1. Verify Workflow Configuration

- [x] 1.1 Verify workflow triggers on push to main
- [x] 1.2 Verify test gate runs before build jobs
- [x] 1.3 Verify cross-compilation for linux/amd64, linux/arm64, linux/arm/v7
- [x] 1.4 Verify artifact naming: foghorn-daemon-linux-<arch>, foghorn-tui-linux-<arch>
- [x] 1.5 Verify Docker images pushed to ghcr.io per architecture
- [x] 1.6 Verify multi-arch manifest with `latest` and SHA tags

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/github-actions-release.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
