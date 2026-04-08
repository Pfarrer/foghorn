## 1. Verify Config Schema

- [ ] 1.1 Verify check config fields: name, image, schedule, evaluation, description, tags, enabled, env, timeout, metadata
- [ ] 1.2 Verify schedule formats: cron and interval
- [ ] 1.3 Verify evaluation rules: type, condition, threshold, expected, metadata
- [ ] 1.4 Verify global config fields parsed correctly
- [ ] 1.5 Verify env map passed to check containers

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/define-config-yaml-format.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
