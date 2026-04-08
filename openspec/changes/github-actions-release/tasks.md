## 1. Verify Workflow Configuration

- [ ] 1.1 Verify workflow triggers on push to main
- [ ] 1.2 Verify test gate runs before build jobs
- [ ] 1.3 Verify cross-compilation for linux/amd64, linux/arm64, linux/arm/v7
- [ ] 1.4 Verify artifact naming: foghorn-daemon-linux-<arch>, foghorn-tui-linux-<arch>
- [ ] 1.5 Verify Docker images pushed to ghcr.io per architecture
- [ ] 1.6 Verify multi-arch manifest with `latest` and SHA tags

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/github-actions-release.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
