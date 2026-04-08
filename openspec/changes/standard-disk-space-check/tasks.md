## 1. Verify Existing Implementation Against Specs

- [ ] 1.1 Verify mount point validation: missing MOUNT_POINT → unknown, non-existent path → unknown
- [ ] 1.2 Verify disk space stats: df output parsed correctly into total_bytes, used_bytes, free_bytes, usage_percent
- [ ] 1.3 Verify percentage thresholds: below warning → pass, at/above warning → warn, at/above critical → fail
- [ ] 1.4 Verify invalid threshold relationship: critical < warning → unknown
- [ ] 1.5 Verify byte-based thresholds: WARNING_THRESHOLD_BYTES and CRITICAL_THRESHOLD_BYTES evaluated alongside percentages
- [ ] 1.6 Verify inode checking: enabled by default, inode stats populated, thresholds enforced
- [ ] 1.7 Verify inode graceful degradation: filesystems reporting `-` for inodes → inode check disabled with note
- [ ] 1.8 Verify CHECK_INODES=false: inode stats reported as zero, message notes disabled
- [ ] 1.9 Verify timeout handling: TIMEOUT_SECONDS default 10, FOGHORN_TIMEOUT override
- [ ] 1.10 Verify JSON output: valid JSON with all required fields (status, message, data, timestamp, duration_ms)
- [ ] 1.11 Verify exit codes: pass → 0, warn/fail/unknown → 1

## 2. Documentation

- [ ] 2.1 Delete legacy spec file `specs/standard-disk-space-check.md`
- [ ] 2.2 Update `specs/STATUS.md` to remove the standard-disk-space-check entry from the Done section
- [ ] 2.3 Verify all Go tests pass and code compiles without warnings

## 3. Archive Change

- [ ] 3.1 Archive the OpenSpec change once verification is complete
