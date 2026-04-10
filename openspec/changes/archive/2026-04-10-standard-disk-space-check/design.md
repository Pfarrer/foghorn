## Context

The disk-check container (`containers/disk-check/`) is a fully implemented Docker container that monitors filesystem disk space and inode usage. It runs on Alpine Linux using `df` for stats. The current implementation matches the legacy spec closely — it supports percentage and byte-based thresholds, optional inode checking, and proper timeout handling.

## Goals / Non-Goals

**Goals:**
- Document the existing disk-check implementation as formal OpenSpec artifacts
- Preserve all current behavior in structured, testable spec requirements
- Enable future changes through the OpenSpec workflow

**Non-Goals:**
- Modifying the container implementation
- Adding multi-mount-point support or new features
- Changing the base image

## Decisions

**Alpine + df over Go implementation**: The current implementation uses `check.sh` on `alpine:latest` with POSIX `df`. Minimal image, no additional dependencies needed.

**Dual threshold model**: Both percentage-based (`WARNING_THRESHOLD_PERCENT`, `CRITICAL_THRESHOLD_PERCENT`) and byte-based (`WARNING_THRESHOLD_BYTES`, `CRITICAL_THRESHOLD_BYTES`) thresholds are supported. Byte thresholds are checked alongside percentage thresholds — the most severe status wins.

**Inode checking**: Enabled by default (`CHECK_INODES=true`). Uses `df -Pi` for inode stats. Gracefully handles filesystems that report `-` for inode values by disabling inode checks with a note.

**Host filesystem access**: Designed to run with `/:/host:ro` volume mount. `MOUNT_POINT` references paths under `/host` (e.g., `/host/var/log`).

## Risks / Trade-offs

- [Migration doc-only] → No functional risk; artifacts describe existing behavior
- [Legacy spec coexists] → Both legacy and OpenSpec specs exist until legacy is retired
