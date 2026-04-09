## 1. Verify Config Schema

- [x] 1.1 Verify check config fields: name, image, schedule, evaluation, description, tags, enabled, env, timeout, metadata
- [x] 1.2 Verify schedule formats: cron and interval
- [x] 1.3 Verify evaluation rules: type, condition, threshold, expected, metadata
- [x] 1.4 Verify global config fields parsed correctly
- [x] 1.5 Verify env map passed to check containers

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/define-config-yaml-format.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
