# Standard System Load Check Container

## Category
integration

## Description
Define a reusable Docker container that monitors system load average and CPU usage

## Usage Steps
1. User pulls the standard system load check image: `foghorn/load-check:latest`
2. Configure check with load thresholds and time intervals
3. Foghorn runs the container on schedule
4. Container reads load metrics and outputs results

## Implementation Notes
Container should accept these environment variables:
- `WARNING_THRESHOLD_1M`: 1-minute load average threshold for "warn" (default: number of CPU cores)
- `CRITICAL_THRESHOLD_1M`: 1-minute load average threshold for "fail" (default: 2x CPU cores)
- `WARNING_THRESHOLD_5M`: 5-minute load average threshold for "warn"
- `CRITICAL_THRESHOLD_5M`: 5-minute load average threshold for "fail"
- `WARNING_THRESHOLD_15M`: 15-minute load average threshold for "warn"
- `CRITICAL_THRESHOLD_15M`: 15-minute load average threshold for "fail"
- `CPU_COUNT`: Number of CPU cores to normalize against (auto-detect if not specified)
- `NORMALIZE_LOAD`: Whether to normalize load by CPU count (default: true)
- `INCLUDE_CPU_PERCENT`: Whether to include CPU usage percentage (default: false)

Note: Container must run with access to host metrics:
```yaml
volumes:
  - /proc:/host/proc:ro
pid: "host"
```

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Load average: 2.15 (1m), 1.98 (5m), 1.45 (15m) on 4 cores (54%)",
  "data": {
    "load_1m": 2.15,
    "load_5m": 1.98,
    "load_15m": 1.45,
    "cpu_count": 4,
    "load_1m_normalized": 0.54,
    "load_5m_normalized": 0.50,
    "load_15m_normalized": 0.36,
    "cpu_usage_percent": 42.5
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 32
}
```

Status determination (when NORMALIZE_LOAD=true):
- `pass`: All normalized loads < WARNING_THRESHOLD (default 1.0 = 100% of CPU capacity)
- `warn`: Any normalized load >= WARNING_THRESHOLD but < CRITICAL_THRESHOLD
- `fail`: Any normalized load >= CRITICAL_THRESHOLD (default 2.0 = 200% of CPU capacity)
- `unknown`: Error reading load information

Status determination (when NORMALIZE_LOAD=false):
- `pass`: All loads < WARNING_THRESHOLD (absolute values)
- `warn`: Any load >= WARNING_THRESHOLD but < CRITICAL_THRESHOLD
- `fail`: Any load >= CRITICAL_THRESHOLD (absolute values)
- `unknown`: Error reading load information

Use Alpine Linux with `/proc/loadavg` parsing or Go for cross-platform

## Acceptance Criteria
- [ ] Container reads load averages from /proc/loadavg
- [ ] Reports all three load averages (1m, 5m, 15m)
- [ ] Auto-detects CPU count if not specified
- [ ] Optionally normalizes load by CPU count
- [ ] Returns appropriate status based on thresholds
- [ ] Includes both absolute and normalized values in data
- [ ] Optionally includes CPU usage percentage
- [ ] Works with host /proc mount
- [ ] Example configurations for different scenarios
- [ ] Respects FOGHORN_TIMEOUT if provided

## Example Configuration
```yaml
name: "system-load-check"
image: "foghorn/load-check:latest"
schedule:
  interval: "2m"
env:
  WARNING_THRESHOLD_1M: "2.0"
  CRITICAL_THRESHOLD_1M: "4.0"
  NORMALIZE_LOAD: "true"
timeout: "10s"
volumes:
  - /proc:/host/proc:ro
pid: "host"
```

```yaml
name: "cpu-intensive-monitor"
image: "foghorn/load-check:latest"
schedule:
  cron: "*/5 * * * *"
env:
  WARNING_THRESHOLD_1M: "3.0"
  CRITICAL_THRESHOLD_1M: "6.0"
  INCLUDE_CPU_PERCENT: "true"
timeout: "10s"
volumes:
  - /proc:/host/proc:ro
pid: "host"
```

## Passes
false
