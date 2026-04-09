## 1. Verify Version Selectors

- [x] 1.1 Verify selector parsing: MAJOR, MAJOR.PATCH, MAJOR.MINOR.PATCH, invalid formats
- [x] 1.2 Verify partial selector resolution to latest matching version
- [x] 1.3 Verify full selector exact matching
- [x] 1.4 Verify no-match returns false
- [x] 1.5 Verify image reference parsing: valid, digest rejected, latest rejected

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/versioned-check-container-releases.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
