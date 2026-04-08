## 1. Verify Timeout Handling

- [ ] 1.1 Verify per-check timeout parsed from config, defaults to 30s when unset
- [ ] 1.2 Verify container terminated when timeout exceeded
- [ ] 1.3 Verify `FOGHORN_TIMEOUT` env var injected with parsed timeout value
- [ ] 1.4 Verify timeout events logged at warn level

## 2. Verify Concurrency

- [ ] 2.1 Verify checks queued when `max_concurrent_checks` limit reached
- [ ] 2.2 Verify queued checks dispatched on slot availability
- [ ] 2.3 Verify no resource exhaustion with many checks

## 3. Documentation

- [ ] 3.1 Delete legacy spec file `specs/concurrent-execution-timeout.md`
- [ ] 3.2 Update `specs/STATUS.md` to remove the entry
- [ ] 3.3 Verify all tests pass and code compiles without warnings

## 4. Archive Change

- [ ] 4.1 Archive the OpenSpec change once verification is complete
