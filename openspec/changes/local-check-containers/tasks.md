## 1. Verify Local Containers

- [ ] 1.1 Verify each container has Dockerfile and README.md
- [ ] 1.2 Verify workflow triggers on `containers/**` path changes
- [ ] 1.3 Verify change detection builds only modified containers
- [ ] 1.4 Verify unchanged containers are not rebuilt
- [ ] 1.5 Verify no-container-changes push skips the build job
- [ ] 1.6 Verify VERSION file read and semver validation
- [ ] 1.7 Verify image tagged with version and `latest`

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/local-check-containers.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
