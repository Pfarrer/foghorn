## 1. Verify Debug Output

- [ ] 1.1 Verify mode behavior: off, on_failure (pass/fail/timeout), always
- [ ] 1.2 Verify per-check override of global mode
- [ ] 1.3 Verify secret values redacted from logged output
- [ ] 1.4 Verify auth patterns redacted from logged output
- [ ] 1.5 Verify output truncation at configured max chars

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/check-container-debug-output.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
