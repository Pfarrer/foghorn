## 1. Verify Existing Implementation Against Specs

- [x] 1.1 Verify mount point validation: missing MOUNT_POINT → unknown, non-existent path → unknown
- [x] 1.2 Verify disk space stats: df output parsed correctly into total_bytes, used_bytes, free_bytes, usage_percent
- [x] 1.3 Verify percentage thresholds: below warning → pass, at/above warning → warn, at/above critical → fail
- [x] 1.4 Verify invalid threshold relationship: critical < warning → unknown
- [x] 1.5 Verify byte-based thresholds: WARNING_THRESHOLD_BYTES and CRITICAL_THRESHOLD_BYTES evaluated alongside percentages
- [x] 1.6 Verify inode checking: enabled by default, inode stats populated, thresholds enforced
- [x] 1.7 Verify inode graceful degradation: filesystems reporting `-` for inodes → inode check disabled with note
- [x] 1.8 Verify CHECK_INODES=false: inode stats reported as zero, message notes disabled
- [x] 1.9 Verify timeout handling: TIMEOUT_SECONDS default 10, FOGHORN_TIMEOUT override
- [x] 1.10 Verify JSON output: valid JSON with all required fields (status, message, data, timestamp, duration_ms)
- [x] 1.11 Verify exit codes: pass → 0, warn/fail/unknown → 1

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/standard-disk-space-check.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the standard-disk-space-check entry from the Done section
- [x] 2.3 Verify all Go tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
