## 1. Verify CLI Flags

- [x] 1.1 Verify `-h`/`--help` prints usage and exits 0
- [x] 1.2 Verify `-c`/`--config` loads config; omitting it prints error and exits 1
- [x] 1.3 Verify `-v`/`--verbose` enables verbose logging
- [x] 1.4 Verify `-d`/`--dry-run` validates config without starting scheduler
- [x] 1.5 Verify `-l`/`--log-level` accepts debug/info/warn/error; rejects invalid values
- [x] 1.6 Verify all flags work in combination

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/basic-cli-arguments.md` (already doesn't exist)
- [x] 2.2 Update `specs/STATUS.md` to remove the entry (already doesn't exist)
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
