## 1. Verify Secret Injection

- [ ] 1.1 Verify `secret://` reference format parsed correctly
- [ ] 1.2 Verify non-secret values pass through unchanged
- [ ] 1.3 Verify secrets encrypted at rest with AES-256-GCM
- [ ] 1.4 Verify secrets resolved only in memory at runtime
- [ ] 1.5 Verify secrets not written to state logs or history
- [ ] 1.6 Verify ephemeral file injection via `_FILE` env vars
- [ ] 1.7 Verify secrets never passed as Docker CLI arguments
- [ ] 1.8 Verify CLI: set, list (keys only), delete

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/secret-injection-for-check-containers.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
