## 1. Verify Config Loading

- [x] 1.1 Verify valid YAML file loads correctly
- [x] 1.2 Verify file-not-found and invalid YAML produce errors
- [x] 1.3 Verify document-per-check format parsing
- [x] 1.4 Verify list format parsing
- [x] 1.5 Verify mixed formats in one file
- [x] 1.6 Verify validation: missing name, missing image, missing schedule, both cron+interval, invalid image tag, negative max_concurrent
- [x] 1.7 Verify config summary printed on success
- [x] 1.8 Verify single-load (no hot-reload)

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/load-config-yaml.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
