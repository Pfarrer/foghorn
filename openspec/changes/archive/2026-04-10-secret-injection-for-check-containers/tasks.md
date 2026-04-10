## 1. Verify Secret Injection

- [x] 1.1 Verify `secret://` reference format parsed correctly
- [x] 1.2 Verify non-secret values pass through unchanged
- [x] 1.3 Verify secrets encrypted at rest with AES-256-GCM
- [x] 1.4 Verify secrets resolved only in memory at runtime
- [x] 1.5 Verify secrets not written to state logs or history
- [x] 1.6 Verify ephemeral file injection via `_FILE` env vars
- [x] 1.7 Verify secrets never passed as Docker CLI arguments
- [x] 1.8 Verify CLI: set, list (keys only), delete

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/secret-injection-for-check-containers.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
