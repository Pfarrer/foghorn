## 1. Verify Interval Scheduling

- [ ] 1.1 Verify interval parsing for s, m, h, d units
- [ ] 1.2 Verify invalid intervals produce errors (negative, empty, bad unit, non-numeric)
- [ ] 1.3 Verify interval checks fire and NextRun advances by interval duration
- [ ] 1.4 Verify immediate first run on empty state
- [ ] 1.5 Verify mixed cron and interval schedules coexist

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/interval-based-scheduling.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
