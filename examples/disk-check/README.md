# Standard Disk Space Check

A Docker container that monitors filesystem free space and usage as part of the Foghorn monitoring system.

## Features

- Monitors disk space usage on any mount point
- Supports both percentage and byte-based thresholds
- Optional inode usage monitoring
- Works with read-only mounted host filesystems
- Respects Foghorn timeout settings

## Environment Variables

### Required
- `MOUNT_POINT`: Filesystem mount point to check (e.g., `/`, `/var/log`). When checking host filesystems, prefix with `/host` (e.g., `/host/var/log`).

### Optional
- `WARNING_THRESHOLD_PERCENT`: Usage percentage for "warn" status (default: 80)
- `CRITICAL_THRESHOLD_PERCENT`: Usage percentage for "fail" status (default: 90)
- `WARNING_THRESHOLD_BYTES`: Usage in bytes for "warn" status (alternative to percent)
- `CRITICAL_THRESHOLD_BYTES`: Usage in bytes for "fail" status (alternative to percent)
- `INODE_WARNING_PERCENT`: Inode usage percentage for "warn" status (default: 85)
- `INODE_CRITICAL_PERCENT`: Inode usage percentage for "fail" status (default: 95)
- `CHECK_INODES`: Whether to check inode usage (default: true)
- `FOGHORN_TIMEOUT`: Timeout duration for the check (from Foghorn)

## Volume Mounts

To check the host filesystem, mount it as read-only:

```yaml
volumes:
  - /:/host:ro
```

And set `MOUNT_POINT` with the `/host` prefix.

**Note**: When checking host filesystems from within a container, inode information may not be available due to Docker's overlay filesystem. In such cases, inode checks will be skipped automatically.

## Status Determination

- **pass**: Usage percent < WARNING_THRESHOLD_PERCENT AND inode usage < INODE_WARNING_PERCENT
- **warn**: Usage percent >= WARNING_THRESHOLD_PERCENT OR inode usage >= INODE_WARNING_PERCENT (but below critical thresholds)
- **fail**: Usage percent >= CRITICAL_THRESHOLD_PERCENT OR inode usage >= INODE_CRITICAL_PERCENT
- **unknown**: Error reading filesystem information or mount point doesn't exist

## Output Format

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

## Building

```bash
cd examples/disk-check
docker build -t foghorn/disk-check:1.0.0 .
```

## Example Configurations

### Check root filesystem

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
volumes:
  - /:/host:ro
```

### Check logs directory with inode monitoring

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

### Check with byte-based thresholds

```yaml
name: "data-disk-space"
image: "foghorn/disk-check:1.0.0"
schedule:
  cron: "0 * * * *"
env:
  MOUNT_POINT: "/host/data"
  WARNING_THRESHOLD_BYTES: "100000000000"
  CRITICAL_THRESHOLD_BYTES: "200000000000"
timeout: "10s"
volumes:
  - /:/host:ro
```

## Testing Locally

```bash
# Check current directory
docker run --rm \
  -v $PWD:/host:ro \
  -e MOUNT_POINT=/host \
  foghorn/disk-check:1.0.0

# Check root filesystem
docker run --rm \
  -v /:/host:ro \
  -e MOUNT_POINT=/host \
  -e WARNING_THRESHOLD_PERCENT=70 \
  -e CRITICAL_THRESHOLD_PERCENT=85 \
  foghorn/disk-check:1.0.0
```