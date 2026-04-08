## 1. Verify Concurrency Limiting

- [ ] 1.1 Verify `max_concurrent_checks` config field: positive value limits concurrency, 0/unset = unlimited
- [ ] 1.2 Verify negative `max_concurrent_checks` rejected at config load
- [ ] 1.3 Verify queued checks dispatch when running slots become available

## 2. Verify Priority Ordering

- [ ] 2.1 Verify due checks sorted by priority duration (longer interval first)
- [ ] 2.2 Verify equal-interval tie-breaker uses alphabetical check name
- [ ] 2.3 Verify queue sorting (`sortQueueLocked`) applies same priority logic

## 3. Verify Initial Run Behavior

- [ ] 3.1 Verify interval checks run immediately on empty state
- [ ] 3.2 Verify cron checks wait for scheduled time on empty state

## 4. Documentation

- [ ] 4.1 Delete legacy spec file `specs/check-execution-order-concurrency.md`
- [ ] 4.2 Update `specs/STATUS.md` to remove the entry
- [ ] 4.3 Verify all tests pass and code compiles without warnings

## 5. Archive Change

- [ ] 5.1 Archive the OpenSpec change once verification is complete
