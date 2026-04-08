## 1. Verify State Log

- [ ] 1.1 Verify config validation: state_log_file requires state_log_period, invalid period rejected
- [ ] 1.2 Verify results persisted as JSON lines
- [ ] 1.3 Verify retention pruning on load and append
- [ ] 1.4 Verify scheduler state restored from latest per-check results on restart
- [ ] 1.5 Verify missing file creates new file and starts fresh
- [ ] 1.6 Verify corrupt file handled gracefully
- [ ] 1.7 Verify file locking prevents concurrent access

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/check-state-log-file.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
