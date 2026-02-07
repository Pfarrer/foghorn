# Standard HTTP Status Check Container

## Category
integration

## Description
Define a reusable Docker container that performs HTTP/HTTPS status checks on web endpoints

## Usage Steps
1. User pulls the standard HTTP check image: `foghorn/http-check:1.0.0`
2. Configure check with target URL and expected status/behavior
3. Foghorn runs the container on schedule
4. Container performs HTTP request and outputs results

## Implementation Notes
Container should accept these environment variables:
- `URL` (required): Target URL to check
- `EXPECTED_STATUS`: Expected HTTP status code or code range (default: 200)
  - Can be single code: `200`
  - Can be range: `200-299`
  - Can be comma-separated: `200,201,204`
- `METHOD`: HTTP method to use (default: GET)
- `HEADERS`: JSON string with request headers (e.g., `{"Authorization":"Bearer token"}`)
- `REQUEST_BODY`: Body content for POST/PUT requests
- `TIMEOUT_SECONDS`: Request timeout (default: 30)
- `FOLLOW_REDIRECTS`: Whether to follow HTTP redirects (default: true)
- `VERIFY_SSL`: Whether to verify SSL certificates (default: true)
- `WARNING_THRESHOLD_MS`: Response time threshold for "warn" status (default: 1000ms)
- `CRITICAL_THRESHOLD_MS`: Response time threshold for "fail" status (default: 5000ms)
- `CONTENT_REGEX`: Optional regex pattern to validate response body

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "HTTP GET https://api.example.com/health: 200 OK (234ms)",
  "data": {
    "url": "https://api.example.com/health",
    "method": "GET",
    "status_code": 200,
    "status_text": "OK",
    "response_time_ms": 234,
    "response_size_bytes": 1024,
    "content_match": true
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 250
}
```

Status determination:
- `pass`: Status code matches expected, response time < WARNING_THRESHOLD_MS, content regex passes if specified
- `warn`: Status code matches expected, response time >= WARNING_THRESHOLD_MS
- `fail`: Status code doesn't match expected, response time >= CRITICAL_THRESHOLD_MS, connection error, or content regex fails
- `unknown`: Error executing request or invalid configuration

Use Alpine Linux with `curl` or Go-based implementation for better control

## Acceptance Criteria
- [ ] Container accepts URL and validates it's a valid HTTP/HTTPS URL
- [ ] Performs HTTP request with configurable method and headers
- [ ] Supports status code comparison (single, range, or list)
- [ ] Measures and reports response time
- [ ] Supports optional content validation with regex
- [ ] Respects SSL verification setting
- [ ] Follows redirects by default with option to disable
- [ ] Returns appropriate status based on thresholds and response code
- [ ] Includes response size in data object
- [ ] Handles connection timeouts gracefully
- [ ] Respects FOGHORN_TIMEOUT if provided
- [ ] Example configurations for common scenarios (health checks, API checks)

## Example Configuration
```yaml
name: "api-health-check"
image: "foghorn/http-check:1.0.0"
schedule:
  interval: "1m"
env:
  URL: "https://api.example.com/health"
  EXPECTED_STATUS: "200-299"
  WARNING_THRESHOLD_MS: "1000"
  CRITICAL_THRESHOLD_MS: "5000"
timeout: "30s"
```

```yaml
name: "website-check"
image: "foghorn/http-check:1.0.0"
schedule:
  cron: "*/10 * * * *"
env:
  URL: "https://example.com"
  CONTENT_REGEX: ".*Welcome.*"
  HEADERS: '{"User-Agent":"Foghorn/1.0"}'
timeout: "15s"
```

## Passes
true
