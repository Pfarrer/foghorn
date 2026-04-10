## 1. Verify Timeout Handling

- [x] 1.1 Verify per-check timeout parsed from config, defaults to 30s when unset
- [x] 1.2 Verify container terminated when timeout exceeded
- [x] 1.3 Verify `FOGHORN_TIMEOUT` env var injected with parsed timeout value
- [x] 1.4 Verify timeout events logged at warn level

## 2. Verify Concurrency

- [x] 2.1 Verify checks queued when `max_concurrent_checks` limit reached
- [x] 2.2 Verify queued checks dispatched on slot availability
- [x] 2.3 Verify no resource exhaustion with many checks

## 3. Documentation

- [x] 3.1 Delete legacy spec file `specs/concurrent-execution-timeout.md`
- [x] 3.2 Update `specs/STATUS.md` to remove the entry
- [x] 3.3 Verify all tests pass and code compiles without warnings

## 4. Archive Change

- [x] 4.1 Archive the OpenSpec change once verification is complete
