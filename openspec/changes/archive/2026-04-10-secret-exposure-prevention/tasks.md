## 1. Create Redact Package

- [x] 1.1 Create `redact/` package with `Sanitize(input string, secrets []string) string` function
- [x] 1.2 Add common pattern matching: auth headers, password=, token=, bearer tokens
- [x] 1.3 Add tests for sanitize function

## 2. Update Logger

- [x] 2.1 Integrate redact package into `logger/logger.go`
- [x] 2.2 Verify no secret values appear in logs at any level

## 3. Update Status API

- [x] 3.1 Review `internal/statusapi/` for any check config fields that may include secrets
- [x] 3.2 Strip secret values from API responses before serialization

## 4. Runtime Guardrails

- [x] 4.1 Add validation: before logging check env vars, verify no raw secret values
- [x] 4.2 If detected, log a warning and use redacted placeholder

## 5. Integration Tests

- [x] 5.1 Add test that runs checks with secrets and inspects log output
- [x] 5.2 Add test that fetches status API and verifies no secrets
- [x] 5.3 Add test that verifies `ps` output contains no secret values

## 6. Verification

- [x] 6.1 Verify code compiles without warnings and all tests pass

## 7. Documentation

- [x] 7.1 Delete legacy spec file `specs/secret-exposure-prevention.md`
- [x] 7.2 Update `specs/STATUS.md` to remove the entry

## 8. Archive Change

- [x] 8.1 Archive the OpenSpec change once implementation is complete
