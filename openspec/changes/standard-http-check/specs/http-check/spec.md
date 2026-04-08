## ADDED Requirements

### Requirement: URL validation
The container SHALL validate that the `URL` environment variable is provided and contains a valid HTTP or HTTPS URL. Missing or invalid URLs SHALL produce a `fail` status.

#### Scenario: Missing URL
- **WHEN** neither `URL` nor `CHECK_URL` environment variable is set
- **THEN** the container outputs JSON with status `fail` and message indicating URL is required

#### Scenario: Invalid URL scheme
- **WHEN** `URL` is set to a value that does not start with `http://` or `https://`
- **THEN** the container outputs JSON with status `fail` and message indicating invalid URL

#### Scenario: Backward compatible CHECK_URL
- **WHEN** `URL` is not set but `CHECK_URL` is set
- **THEN** the container uses `CHECK_URL` as the target URL

### Requirement: Configurable HTTP method
The container SHALL support the `METHOD` environment variable to set the HTTP method (default: `GET`). The method SHALL be passed directly to curl.

#### Scenario: Custom method
- **WHEN** `METHOD` is set to `POST`
- **THEN** the container sends an HTTP POST request to the target URL

#### Scenario: Default method
- **WHEN** `METHOD` is not set
- **THEN** the container sends an HTTP GET request

### Requirement: Custom request headers
The container SHALL support the `HEADERS` environment variable as a JSON string of request headers. Each key-value pair SHALL be passed to curl as an `-H` flag.

#### Scenario: Custom headers
- **WHEN** `HEADERS` is set to `{"Authorization":"Bearer token","Accept":"application/json"}`
- **THEN** both headers are included in the HTTP request

#### Scenario: No headers
- **WHEN** `HEADERS` is not set
- **THEN** no extra headers are added beyond curl defaults

### Requirement: Request body
The container SHALL support the `REQUEST_BODY` environment variable for POST/PUT request bodies, passed to curl via `-d`.

#### Scenario: POST with body
- **WHEN** `METHOD` is `POST` and `REQUEST_BODY` is `{"key":"value"}`
- **THEN** the container sends the body with the request

#### Scenario: No body
- **WHEN** `REQUEST_BODY` is not set
- **THEN** no request body is sent

### Requirement: Status code matching
The container SHALL support `EXPECTED_STATUS` in three formats: single code (`200`), range (`200-299`), and comma-separated list (`200,201,204`). Default SHALL be `200`. Status codes not matching any format SHALL produce `fail`.

#### Scenario: Single status code match
- **WHEN** `EXPECTED_STATUS` is `200` and the server responds with 200
- **THEN** status code check passes

#### Scenario: Single status code mismatch
- **WHEN** `EXPECTED_STATUS` is `200` and the server responds with 404
- **THEN** the container outputs status `fail` with message indicating the mismatch

#### Scenario: Status code range match
- **WHEN** `EXPECTED_STATUS` is `200-299` and the server responds with 204
- **THEN** status code check passes

#### Scenario: Status code range mismatch
- **WHEN** `EXPECTED_STATUS` is `200-299` and the server responds with 301
- **THEN** the container outputs status `fail`

#### Scenario: Comma-separated status codes
- **WHEN** `EXPECTED_STATUS` is `200,201,204` and the server responds with 201
- **THEN** status code check passes

### Requirement: Response time thresholds
The container SHALL measure response time and compare against `WARNING_THRESHOLD_MS` (default 1000) and `CRITICAL_THRESHOLD_MS` (default 5000). When status code matches: response time below warning threshold → `pass`; at or above warning threshold but below critical → `warn`; at or above critical threshold → `fail`.

#### Scenario: Response under warning threshold
- **WHEN** status code matches and response time is 500ms (below 1000ms warning threshold)
- **THEN** container outputs status `pass`

#### Scenario: Response at or above warning threshold
- **WHEN** status code matches and response time is 1200ms (above 1000ms warning, below 5000ms critical)
- **THEN** container outputs status `warn` with message indicating slow response

#### Scenario: Response at or above critical threshold
- **WHEN** status code matches and response time is 5500ms (above 5000ms critical threshold)
- **THEN** container outputs status `fail` with message indicating critical response time

### Requirement: Content regex validation
The container SHALL support optional `CONTENT_REGEX` to validate the response body. When set, the response body SHALL be matched against the regex. A match SHALL set `content_match: true`; no match SHALL produce status `fail` with `content_match: false`.

#### Scenario: Regex matches
- **WHEN** `CONTENT_REGEX` is `.*Welcome.*` and the response body contains "Welcome"
- **THEN** `content_match` is `true` in the output data

#### Scenario: Regex does not match
- **WHEN** `CONTENT_REGEX` is `.*Welcome.*` and the response body does not contain "Welcome"
- **THEN** the container outputs status `fail` and `content_match` is `false`

#### Scenario: No regex configured
- **WHEN** `CONTENT_REGEX` is not set
- **THEN** `content_match` is omitted from output or set to `null`

### Requirement: SSL verification control
The container SHALL support `VERIFY_SSL` (default: `true`). When `false`, the container SHALL skip TLS certificate verification (curl `-k` flag).

#### Scenario: SSL verification enabled
- **WHEN** `VERIFY_SSL` is `true` (default)
- **THEN** the container verifies SSL certificates and fails on invalid certs

#### Scenario: SSL verification disabled
- **WHEN** `VERIFY_SSL` is `false`
- **THEN** the container skips SSL certificate verification

### Requirement: Redirect control
The container SHALL support `FOLLOW_REDIRECTS` (default: `true`). When `true`, curl follows redirects (`-L` flag). When `false`, redirects are not followed.

#### Scenario: Follow redirects enabled
- **WHEN** `FOLLOW_REDIRECTS` is `true` and the server responds with 301
- **THEN** curl follows the redirect and reports the final response

#### Scenario: Follow redirects disabled
- **WHEN** `FOLLOW_REDIRECTS` is `false` and the server responds with 301
- **THEN** the container reports the 301 status code without following

### Requirement: Timeout handling
The container SHALL respect `TIMEOUT_SECONDS` (default: 30) and `FOGHORN_TIMEOUT` for request timeouts. If `FOGHORN_TIMEOUT` is set and is shorter than `TIMEOUT_SECONDS`, it SHALL take precedence. Timeouts SHALL produce a `fail` status.

#### Scenario: Connection timeout
- **WHEN** the HTTP request does not complete within the configured timeout
- **THEN** the container outputs JSON with status `fail` and message indicating timeout

#### Scenario: FOGHORN_TIMEOUT overrides
- **WHEN** `FOGHORN_TIMEOUT` is set to a value shorter than `TIMEOUT_SECONDS`
- **THEN** the effective timeout is `FOGHORN_TIMEOUT`

### Requirement: Standard JSON output format
The container SHALL output JSON to stdout conforming to the Foghorn check contract: `{status, message, data, timestamp, duration_ms}`. The `data` object SHALL include `url`, `method`, `status_code`, `status_text`, `response_time_ms`, `response_size_bytes`, and `content_match` (when regex is configured).

#### Scenario: Valid JSON output on success
- **WHEN** the check completes successfully
- **THEN** stdout is valid JSON with all required fields populated

#### Scenario: Valid JSON output on failure
- **WHEN** the check encounters an error
- **THEN** stdout is valid JSON with status `fail` or `unknown` and a descriptive message

### Requirement: Response data reporting
The container SHALL report the HTTP status code, status text, response time in milliseconds, and response size in bytes in the output `data` object.

#### Scenario: Successful request data
- **WHEN** the HTTP request returns 200 OK with a 1024-byte body in 234ms
- **THEN** the data object includes `status_code: 200`, `status_text: "OK"`, `response_time_ms: 234`, `response_size_bytes: 1024`
