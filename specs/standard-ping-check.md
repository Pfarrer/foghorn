# Standard Ping Check Container

## Category
integration

## Description
Define a reusable Docker container that performs ICMP ping checks to monitor host connectivity and latency

## Usage Steps
1. User pulls the standard ping check image: `foghorn/ping-check:latest`
2. Configure check with target host and optional parameters
3. Foghorn runs the container on schedule
4. Container performs ping and outputs results

## Implementation Notes
Container should accept these environment variables:
- `TARGET_HOST` (required): Hostname or IP address to ping
- `PACKET_COUNT`: Number of ping packets to send (default: 4)
- `PACKET_SIZE`: Size of packets in bytes (default: 56)
- `TIMEOUT_SECONDS`: Total timeout for ping operation (default: 10)
- `WARNING_THRESHOLD_MS`: Latency threshold for "warn" status (default: 100ms)
- `CRITICAL_THRESHOLD_MS`: Latency threshold for "fail" status (default: 500ms)

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Pinged example.com: 4/4 packets received, avg 45ms",
  "data": {
    "host": "example.com",
    "packets_sent": 4,
    "packets_received": 4,
    "packet_loss_percent": 0,
    "min_latency_ms": 42,
    "max_latency_ms": 48,
    "avg_latency_ms": 45
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 1234
}
```

Status determination:
- `pass`: All packets received, avg latency < WARNING_THRESHOLD_MS
- `warn`: All packets received, avg latency >= WARNING_THRESHOLD_MS
- `fail`: Any packet loss OR avg latency >= CRITICAL_THRESHOLD_MS OR host unreachable
- `unknown`: Error executing ping command

Use lightweight Alpine Linux image with `iputils-ping` package

## Acceptance Criteria
- [ ] Container accepts TARGET_HOST and validates it's not empty
- [ ] Performs ping with configurable packet count and timeout
- [ ] Calculates packet loss percentage and latency statistics
- [ ] Returns appropriate status based on thresholds
- [ ] Includes detailed metrics in data object
- [ ] Handles unreachable hosts gracefully with "fail" status
- [ ] Handles DNS resolution failures gracefully
- [ ] Respects FOGHORN_TIMEOUT if provided
- [ ] Works with IPv4 and IPv6 addresses
- [ ] Example configuration in README

## Example Configuration
```yaml
name: "google-ping"
image: "foghorn/ping-check:latest"
schedule:
  cron: "*/5 * * * *"
env:
  TARGET_HOST: "google.com"
  PACKET_COUNT: "4"
  WARNING_THRESHOLD_MS: "100"
  CRITICAL_THRESHOLD_MS: "500"
timeout: "10s"
```

## Passes
false
