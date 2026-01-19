# Standard HTTP Status Check

A Docker container that performs HTTP/HTTPS status checks on web endpoints as part of the Foghorn monitoring system.

## Features

- Supports HTTP and HTTPS endpoints
- Configurable HTTP methods (GET, POST, PUT, DELETE, etc.)
- Status code validation (single code, ranges, or lists)
- Request header customization
- Optional request body for POST/PUT
- Response time measurement with configurable thresholds
- Content validation using regex patterns
- SSL certificate verification control
- HTTP redirect following control
- Respects Foghorn timeout settings

## Environment Variables

### Required
- `URL`: Target HTTP/HTTPS URL to check

### Optional
- `EXPECTED_STATUS`: Expected HTTP status code or code range (default: 200)
  - Single code: `200`
  - Range: `200-299`
  - Comma-separated: `200,201,204`
- `METHOD`: HTTP method to use (default: GET)
- `HEADERS`: JSON string with request headers (e.g., `{"Authorization":"Bearer token"}`)
- `REQUEST_BODY`: Body content for POST/PUT requests
- `TIMEOUT_SECONDS`: Request timeout in seconds (default: 30)
- `FOLLOW_REDIRECTS`: Whether to follow HTTP redirects (default: true)
- `VERIFY_SSL`: Whether to verify SSL certificates (default: true)
- `WARNING_THRESHOLD_MS`: Response time threshold in milliseconds for "warn" status (default: 1000)
- `CRITICAL_THRESHOLD_MS`: Response time threshold in milliseconds for "fail" status (default: 5000)
- `CONTENT_REGEX`: Optional regex pattern to validate response body

## Status Determination

- **pass**: Status code matches expected, response time < WARNING_THRESHOLD_MS, content regex passes if specified
- **warn**: Status code matches expected, response time >= WARNING_THRESHOLD_MS but < CRITICAL_THRESHOLD_MS
- **fail**: Status code doesn't match expected, response time >= CRITICAL_THRESHOLD_MS, connection error, or content regex fails
- **unknown**: Error executing request or invalid configuration

## Output Format

```json
{
  "status": "pass|fail|warn|unknown",
  "message": "HTTP GET https://api.example.com/health: 200 OK (234ms)",
  "data": {
    "url": "https://api.example.com/health",
    "method": "GET",
    "status_code": 200,
    "status_text": "200",
    "response_time_ms": 234,
    "response_size_bytes": 1024,
    "content_match": true
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 250
}
```

## Building

```bash
cd examples/http-check
docker build -t foghorn/http-check:latest .
```

## Example Configurations

### API health check

```yaml
name: "api-health-check"
image: "foghorn/http-check:latest"
schedule:
  interval: "1m"
env:
  URL: "https://api.example.com/health"
  EXPECTED_STATUS: "200-299"
  WARNING_THRESHOLD_MS: "1000"
  CRITICAL_THRESHOLD_MS: "5000"
timeout: "30s"
```

### Website content check

```yaml
name: "website-check"
image: "foghorn/http-check:latest"
schedule:
  cron: "*/10 * * * *"
env:
  URL: "https://example.com"
  CONTENT_REGEX: ".*Welcome.*"
  HEADERS: '{"User-Agent":"Foghorn/1.0"}'
timeout: "15s"
```

### API endpoint with authentication

```yaml
name: "authenticated-api-check"
image: "foghorn/http-check:latest"
schedule:
  interval: "5m"
env:
  URL: "https://api.example.com/users"
  METHOD: "GET"
  EXPECTED_STATUS: "200"
  HEADERS: '{"Authorization":"Bearer YOUR_TOKEN","Accept":"application/json"}'
  WARNING_THRESHOLD_MS: "500"
  CRITICAL_THRESHOLD_MS: "2000"
timeout: "20s"
```

### POST request check

```yaml
name: "api-post-check"
image: "foghorn/http-check:latest"
schedule:
  cron: "0 * * * *"
env:
  URL: "https://api.example.com/webhook"
  METHOD: "POST"
  EXPECTED_STATUS: "200-299"
  REQUEST_BODY: '{"test": true}'
  HEADERS: '{"Content-Type":"application/json"}'
  WARNING_THRESHOLD_MS: "1000"
  CRITICAL_THRESHOLD_MS: "5000"
timeout: "30s"
```

### API with strict timeout

```yaml
name: "fast-api-check"
image: "foghorn/http-check:latest"
schedule:
  interval: "30s"
env:
  URL: "https://api.example.com/status"
  EXPECTED_STATUS: "200"
  TIMEOUT_SECONDS: "5"
  WARNING_THRESHOLD_MS: "100"
  CRITICAL_THRESHOLD_MS: "500"
timeout: "10s"
```

## Testing Locally

```bash
# Simple GET request
docker run --rm \
  -e URL=https://httpbin.org/status/200 \
  foghorn/http-check:latest

# With content validation
docker run --rm \
  -e URL=https://example.com \
  -e CONTENT_REGEX=".*Example.*" \
  foghorn/http-check:latest

# POST request with custom headers
docker run --rm \
  -e URL=https://httpbin.org/post \
  -e METHOD=POST \
  -e REQUEST_BODY='{"test": true}' \
  -e HEADERS='{"Content-Type":"application/json"}' \
  foghorn/http-check:latest

# With status code range
docker run --rm \
  -e URL=https://httpbin.org/status/200 \
  -e EXPECTED_STATUS="200-299" \
  foghorn/http-check:latest
```

## Notes

- The check uses `curl` with SSL verification enabled by default for security
- For development or testing with self-signed certificates, set `VERIFY_SSL=false`
- Response time is measured from when the request starts to when the response is complete
- Content regex is applied to the entire response body
- The check fails immediately on connection errors (DNS resolution failure, refused connections, etc.)
- HTTP redirects are followed by default unless `FOLLOW_REDIRECTS=false`
