## 1. Verify Config Loading

- [ ] 1.1 Verify valid YAML file loads correctly
- [ ] 1.2 Verify file-not-found and invalid YAML produce errors
- [ ] 1.3 Verify document-per-check format parsing
- [ ] 1.4 Verify list format parsing
- [ ] 1.5 Verify mixed formats in one file
- [ ] 1.6 Verify validation: missing name, missing image, missing schedule, both cron+interval, invalid image tag, negative max_concurrent
- [ ] 1.7 Verify config summary printed on success
- [ ] 1.8 Verify single-load (no hot-reload)

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/load-config-yaml.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
