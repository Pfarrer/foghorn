## 1. Verify Image Availability Flag

- [x] 1.1 Verify `-i` and `--verify-image-availability` flags registered
- [x] 1.2 Verify flag combinable with `--dry-run` and other flags
- [x] 1.3 Verify only enabled checks are verified
- [x] 1.4 Verify missing images reported with check names and pull commands
- [x] 1.5 Verify multiple checks using same image are grouped
- [x] 1.6 Verify success message when all images present
- [x] 1.7 Verify Docker daemon connection error handled

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/verify-image-availability-flag.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
