## 1. Verify CLI Flags

- [ ] 1.1 Verify `-h`/`--help` prints usage and exits 0
- [ ] 1.2 Verify `-c`/`--config` loads config; omitting it prints error and exits 1
- [ ] 1.3 Verify `-v`/`--verbose` enables verbose logging
- [ ] 1.4 Verify `-d`/`--dry-run` validates config without starting scheduler
- [ ] 1.5 Verify `-l`/`--log-level` accepts debug/info/warn/error; rejects invalid values
- [ ] 1.6 Verify all flags work in combination

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/basic-cli-arguments.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the entry
- [ ] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
