# Standard Port Check Container

## Category
integration

## Description
Define a reusable Docker container that checks TCP/UDP port connectivity and service availability

## Usage Steps
1. User pulls the standard port check image: `foghorn/port-check:latest`
2. Configure check with target host, port, and protocol
3. Foghorn runs the container on schedule
4. Container attempts connection and outputs results

## Implementation Notes
Container should accept these environment variables:
- `TARGET_HOST` (required): Hostname or IP address to check
- `PORT` (required): Port number to check (1-65535)
- `PROTOCOL`: Protocol to use - `tcp` or `udp` (default: tcp)
- `TIMEOUT_SECONDS`: Connection timeout (default: 5)
- `SEND_DATA`: Data to send after connection (optional, for validation)
- `EXPECT_DATA`: Expected response data (optional, for validation)
- `WARNING_THRESHOLD_MS`: Connection time threshold for "warn" status (default: 1000ms)
- `CRITICAL_THRESHOLD_MS`: Connection time threshold for "fail" status (default: 5000ms)

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "TCP connection to example.com:443 succeeded (45ms)",
  "data": {
    "host": "example.com",
    "port": 443,
    "protocol": "tcp",
    "connected": true,
    "connection_time_ms": 45,
    "bytes_sent": 0,
    "bytes_received": 0,
    "response_data": ""
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 52
}
```

Status determination:
- `pass`: Connection successful, connection time < WARNING_THRESHOLD_MS, response matches EXPECTED_DATA if specified
- `warn`: Connection successful, connection time >= WARNING_THRESHOLD_MS but < CRITICAL_THRESHOLD_MS
- `fail`: Connection failed, connection time >= CRITICAL_THRESHOLD_MS, or response doesn't match expected data
- `unknown`: Error executing check or invalid configuration

For UDP checks:
- Send data if SEND_DATA provided
- Wait for response for specified timeout
- Mark as "pass" if connection possible (UDP is connectionless, so no actual "connection")

Use Alpine Linux with `netcat-openbsd` or Go implementation for better control

## Acceptance Criteria
- [ ] Container accepts TARGET_HOST and PORT, validates port range
- [ ] Supports both TCP and UDP protocols
- [ ] Measures and reports connection time
- [ ] Returns appropriate status based on connection success and thresholds
- [ ] Optionally sends and validates data payload
- [ ] Handles connection timeouts gracefully
- [ ] Reports bytes sent/received when data exchange occurs
- [ ] Works with hostnames and IP addresses
- [ ] Example configurations for common services (SSH, HTTP, database ports)
- [ ] Respects FOGHORN_TIMEOUT if provided

## Example Configuration
```yaml
name: "ssh-port-check"
image: "foghorn/port-check:latest"
schedule:
  cron: "*/2 * * * *"
env:
  TARGET_HOST: "server.example.com"
  PORT: "22"
  PROTOCOL: "tcp"
  WARNING_THRESHOLD_MS: "500"
  CRITICAL_THRESHOLD_MS: "2000"
timeout: "10s"
```

```yaml
name: "database-port-check"
image: "foghorn/port-check:latest"
schedule:
  interval: "30s"
env:
  TARGET_HOST: "db.example.com"
  PORT: "5432"
  PROTOCOL: "tcp"
  SEND_DATA: "\x00\x00\x00\x08\x04\xd2\x16\x2f"
  EXPECT_DATA: "R"
timeout: "5s"
```

## Passes
false
