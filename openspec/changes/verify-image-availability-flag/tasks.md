## 1. Verify Image Availability Flag

- [ ] 1.1 Verify `-i` and `--verify-image-availability` flags registered
- [ ] 1.2 Verify flag combinable with `--dry-run` and other flags
- [ ] 1.3 Verify only enabled checks are verified
- [ ] 1.4 Verify missing images reported with check names and pull commands
- [ ] 1.5 Verify multiple checks using same image are grouped
- [ ] 1.6 Verify success message when all images present
- [ ] 1.7 Verify Docker daemon connection error handled

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/verify-image-availability-flag.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
