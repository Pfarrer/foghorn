## 1. Create Redact Package

- [ ] 1.1 Create `redact/` package with `Sanitize(input string, secrets []string) string` function
- [ ] 1.2 Add common pattern matching: auth headers, password=, token=, bearer tokens
- [ ] 1.3 Add tests for sanitize function

## 2. Update Logger

- [ ] 2.1 Integrate redact package into `logger/logger.go`
- [ ] 2.2 Verify no secret values appear in logs at any level

## 3. Update Status API

- [ ] 3.1 Review `internal/statusapi/` for any check config fields that may include secrets
- [ ] 3.2 Strip secret values from API responses before serialization

## 4. Runtime Guardrails

- [ ] 4.1 Add validation: before logging check env vars, verify no raw secret values
- [ ] 4.2 If detected, log a warning and use redacted placeholder

## 5. Integration Tests

- [ ] 5.1 Add test that runs checks with secrets and inspects log output
- [ ] 5.2 Add test that fetches status API and verifies no secrets
- [ ] 5.3 Add test that verifies `ps` output contains no secret values

## 6. Verification

- [ ] 6.1 Verify code compiles without warnings and all tests pass

## 7. Documentation

- [ ] 7.1 Delete legacy spec file `specs/secret-exposure-prevention.md`
- [ ] 7.2 Update `specs/STATUS.md` to remove the entry

## 8. Archive Change

- [ ] 8.1 Archive the OpenSpec change once implementation is complete
