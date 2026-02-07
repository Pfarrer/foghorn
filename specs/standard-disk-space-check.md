# Standard Disk Space Check Container

## Category
integration

## Description
Define a reusable Docker container that monitors filesystem free space and usage

## Usage Steps
1. User pulls the standard disk space check image: `foghorn/disk-check:1.0.0`
2. Configure check with mount point(s) and thresholds
3. Foghorn runs the container on schedule
4. Container analyzes disk usage and outputs results

## Implementation Notes
Container should accept these environment variables:
- `MOUNT_POINT` (required): Filesystem mount point to check (e.g., `/`, `/var/log`)
- `WARNING_THRESHOLD_PERCENT`: Usage percentage for "warn" status (default: 80)
- `CRITICAL_THRESHOLD_PERCENT`: Usage percentage for "fail" status (default: 90)
- `WARNING_THRESHOLD_BYTES`: Usage in bytes for "warn" status (alternative to percent)
- `CRITICAL_THRESHOLD_BYTES`: Usage in bytes for "fail" status (alternative to percent)
- `INODE_WARNING_PERCENT`: Inode usage percentage for "warn" status (default: 85)
- `INODE_CRITICAL_PERCENT`: Inode usage percentage for "fail" status (default: 95)
- `CHECK_INODES`: Whether to check inode usage (default: true)

Note: Container must run with volume mount to check host filesystem:
```yaml
volumes:
  - /:/host:ro
```
And use `MOUNT_POINT` prefixed with `/host` (e.g., `/host/var/log`)

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "/var/log: 75% used (150GB free of 600GB total), 60% inodes used",
  "data": {
    "mount_point": "/var/log",
    "total_bytes": 644245094400,
    "used_bytes": 483183820800,
    "free_bytes": 161061273600,
    "usage_percent": 75,
    "total_inodes": 10000000,
    "used_inodes": 6000000,
    "free_inodes": 4000000,
    "inode_usage_percent": 60
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 45
}
```

Status determination:
- `pass`: Usage percent < WARNING_THRESHOLD_PERCENT AND inode usage < INODE_WARNING_PERCENT
- `warn`: Usage percent >= WARNING_THRESHOLD_PERCENT OR inode usage >= INODE_WARNING_PERCENT (but below critical thresholds)
- `fail`: Usage percent >= CRITICAL_THRESHOLD_PERCENT OR inode usage >= INODE_CRITICAL_PERCENT
- `unknown`: Error reading filesystem information or mount point doesn't exist

Use Alpine Linux with `df` command or Go for better parsing

## Acceptance Criteria
- [ ] Container accepts MOUNT_POINT and validates it exists
- [ ] Reads disk space statistics (total, used, free)
- [ ] Calculates usage percentage correctly
- [ ] Supports both percentage and byte-based thresholds
- [ ] Optionally checks inode usage
- [ ] Returns appropriate status based on thresholds
- [ ] Includes detailed metrics in data object (bytes and inodes)
- [ ] Handles non-existent mount points gracefully with "unknown" or "fail" status
- [ ] Works with read-only mounted host filesystems
- [ ] Example configurations for different mount points
- [ ] Respects FOGHORN_TIMEOUT if provided

## Example Configuration
```yaml
name: "root-disk-space"
image: "foghorn/disk-check:1.0.0"
schedule:
  cron: "*/10 * * * *"
env:
  MOUNT_POINT: "/host"
  WARNING_THRESHOLD_PERCENT: "80"
  CRITICAL_THRESHOLD_PERCENT: "90"
timeout: "10s"
# Run with host mount
volumes:
  - /:/host:ro
```

```yaml
name: "logs-disk-space"
image: "foghorn/disk-check:1.0.0"
schedule:
  interval: "5m"
env:
  MOUNT_POINT: "/host/var/log"
  WARNING_THRESHOLD_PERCENT: "85"
  CRITICAL_THRESHOLD_PERCENT: "95"
  CHECK_INODES: "true"
  INODE_WARNING_PERCENT: "80"
  INODE_CRITICAL_PERCENT: "90"
timeout: "10s"
volumes:
  - /:/host:ro
```

## Passes
true
