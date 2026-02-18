# Disk Check

Monitors disk usage for a mount point and returns Foghorn JSON output.

## Image
`ghcr.io/pfarrer/foghorn-disk-check:1.0.0`

## Changelog
- 1.0.0 Initial release

## Env
- `MOUNT_POINT` (required)
- `WARNING_THRESHOLD_PERCENT` (default `80`)
- `CRITICAL_THRESHOLD_PERCENT` (default `90`)
- `WARNING_THRESHOLD_BYTES` (optional)
- `CRITICAL_THRESHOLD_BYTES` (optional)
- `CHECK_INODES` (default `true`)
- `INODE_WARNING_PERCENT` (default `85`)
- `INODE_CRITICAL_PERCENT` (default `95`)
- `TIMEOUT_SECONDS` (default `10`)
- `FOGHORN_TIMEOUT` (optional; caps timeout)
