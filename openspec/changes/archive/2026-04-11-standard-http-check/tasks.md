## 1. Input Validation and Configuration

- [x] 1.1 Add URL validation: require `URL` (or `CHECK_URL` fallback), verify http/https scheme, output fail on invalid input
- [x] 1.2 Add `METHOD` env var parsing (default: GET) and pass to curl via `-X`
- [x] 1.3 Add `HEADERS` env var parsing: parse JSON string, generate curl `-H` flags for each key-value pair
- [x] 1.4 Add `REQUEST_BODY` env var support: pass to curl via `-d` when set
- [x] 1.5 Add `TIMEOUT_SECONDS` parsing (default: 30) and `FOGHORN_TIMEOUT` override with min-of logic
- [x] 1.6 Add `VERIFY_SSL` parsing: default true, when false add curl `-k` flag
- [x] 1.7 Add `FOLLOW_REDIRECTS` parsing: default true, when true add curl `-L` flag

## 2. Status Code Matching

- [x] 2.1 Implement `EXPECTED_STATUS` parsing: single code, range (e.g. `200-299`), comma-separated list (e.g. `200,201,204`)
- [x] 2.2 Add status code comparison logic: match actual HTTP code against parsed expected values
- [x] 2.3 Test status matching: single match, range match, list match, mismatch each produces correct result

## 3. Response Time Thresholds

- [x] 3.1 Capture response time from curl write-out (`time_total`), convert to milliseconds
- [x] 3.2 Add `WARNING_THRESHOLD_MS` (default: 1000) and `CRITICAL_THRESHOLD_MS` (default: 5000) env var parsing
- [x] 3.3 Implement threshold logic: pass (< warning), warn (>= warning, < critical), fail (>= critical)
- [x] 3.4 Apply thresholds only when status code matches expected; status mismatch always fails regardless of timing

## 4. Content Regex Validation

- [x] 4.1 Save response body to temp file using curl `-o`
- [x] 4.2 When `CONTENT_REGEX` is set, match body against regex using grep; set `content_match` true/false
- [x] 4.3 When content regex fails to match, output status `fail` with `content_match: false`
- [x] 4.4 When `CONTENT_REGEX` is not set, omit `content_match` from data output

## 5. JSON Output and Response Data

- [x] 5.1 Restructure JSON output to include all data fields: `url`, `method`, `status_code`, `status_text`, `response_time_ms`, `response_size_bytes`, `content_match`
- [x] 5.2 Capture response size from curl write-out (`size_download`)
- [x] 5.3 Capture HTTP status text from response headers
- [x] 5.4 Ensure all JSON output is properly escaped for special characters

## 6. Container Build and Verification

- [x] 6.1 Update Dockerfile if additional packages are needed
- [x] 6.2 Build container and verify it runs: `docker build -t foghorn/http-check:local containers/http-check/`
- [x] 6.3 Test basic GET: `docker run --rm -e URL=https://example.com foghorn/http-check:local`
- [x] 6.4 Test status range: `docker run --rm -e URL=https://example.com -e EXPECTED_STATUS=200-299 foghorn/http-check:local`
- [x] 6.5 Test content regex: `docker run --rm -e URL=https://example.com -e CONTENT_REGEX=".*Example.*" foghorn/http-check:local`
- [x] 6.6 Test timeout handling: `docker run --rm -e URL=https://example.com -e TIMEOUT_SECONDS=1 -e FOGHORN_TIMEOUT=1 foghorn/http-check:local`
- [x] 6.7 Test SSL skip: `docker run --rm -e URL=https://self-signed.badssl.com -e VERIFY_SSL=false foghorn/http-check:local`

## 7. Documentation

- [x] 7.1 Update `specs/STATUS.md` to reference OpenSpec artifacts for standard-http-check
- [x] 7.2 Verify all Go tests pass and code compiles without warnings
