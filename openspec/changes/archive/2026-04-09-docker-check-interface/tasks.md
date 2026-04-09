## 1. Verify Interface Contract

- [x] 1.1 Verify env var injection: FOGHORN_CHECK_NAME, FOGHORN_CHECK_CONFIG, FOGHORN_ENDPOINT, FOGHORN_TIMEOUT
- [x] 1.2 Verify optional env vars omitted when unset
- [x] 1.3 Verify JSON output parsing: valid JSON, invalid JSON, missing fields
- [x] 1.4 Verify exit code semantics: non-zero = failure, zero with pass status = pass
- [x] 1.5 Verify secret injection via `_FILE` env vars pointing to ephemeral files

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/docker-check-interface.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
