## 1. Verify Cron Scheduling

- [ ] 1.1 Verify valid cron expressions parse correctly and set next run time
- [ ] 1.2 Verify invalid cron expressions produce errors on AddCheck
- [ ] 1.3 Verify next run time is in the future and matches cron expression
- [ ] 1.4 Verify checks fire when ticker reaches NextRun
- [ ] 1.5 Verify NextRun is recalculated after execution
- [ ] 1.6 Verify time zone is respected in all time comparisons

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/basic-scheduler-cron.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
