## ADDED Requirements

### Requirement: URL configuration
The container SHALL accept a `URL` environment variable containing the target web page URL. The URL MUST be a valid HTTP or HTTPS URL. Missing or invalid URLs SHALL produce a `fail` status.

#### Scenario: Valid URL provided
- **WHEN** `URL` is set to `https://example.com/page`
- **THEN** the container navigates to that URL and extracts content

#### Scenario: Missing URL
- **WHEN** `URL` is not set
- **THEN** the container outputs JSON with status `fail` and message "URL environment variable is required"

#### Scenario: Invalid URL scheme
- **WHEN** `URL` is set to `ftp://example.com`
- **THEN** the container outputs JSON with status `fail` and message indicating invalid URL

### Requirement: Headless browser rendering
The container SHALL use a headless Chromium browser (via Playwright) to render the page, executing all JavaScript before extracting content.

#### Scenario: JavaScript-rendered content captured
- **WHEN** the target page renders content dynamically via JavaScript
- **THEN** the extracted content includes the dynamically rendered text

#### Scenario: Page load failure
- **WHEN** the browser fails to load the page (DNS error, connection refused, etc.)
- **THEN** the container outputs JSON with status `unknown` and message describing the error (see Technical error handling requirement)

### Requirement: XPath content extraction
The container SHALL accept an optional `XPATH` environment variable. When set, the container SHALL extract only the text content of nodes matching the XPath expression. When unset, the container SHALL extract the full page body text.

#### Scenario: XPath extracts specific content
- **WHEN** `XPATH` is set to `//div[@id='main-content']` and the page contains `<div id="main-content">Hello World</div>`
- **THEN** the extracted content is "Hello World"

#### Scenario: XPath matches multiple nodes
- **WHEN** `XPATH` is set to `//li` and the page has multiple `<li>` elements
- **THEN** the extracted content combines text from all matching nodes, joined by newlines

#### Scenario: XPath with no matches
- **WHEN** `XPATH` is set to `//div[@id='nonexistent']` and no elements match
- **THEN** the container outputs JSON with status `fail` and message indicating no XPath matches

#### Scenario: No XPath configured
- **WHEN** `XPATH` is not set
- **THEN** the container extracts the full text content of the page body

### Requirement: Wait strategy
The container SHALL support `WAIT_SECONDS` (integer, default 0) to wait after page load before extraction, and `WAIT_FOR_SELECTOR` (CSS selector string) to wait until a matching element appears. Both MAY be used together.

#### Scenario: Time-based wait
- **WHEN** `WAIT_SECONDS` is set to `3`
- **THEN** the container waits 3 seconds after page load before extracting content

#### Scenario: Selector-based wait
- **WHEN** `WAIT_FOR_SELECTOR` is set to `.dynamic-content` and the element appears after 2 seconds
- **THEN** the container waits until the element is present, then extracts content

#### Scenario: Default no wait
- **WHEN** neither `WAIT_SECONDS` nor `WAIT_FOR_SELECTOR` is set
- **THEN** content is extracted immediately after page load

### Requirement: State persistence
The container SHALL store the SHA-256 hash of the extracted content in a file at `FOGHORN_PERSISTENT_DIR/state.json`. The file SHALL contain a JSON object with fields `hash`, `timestamp`, and `url`.

#### Scenario: State file created on first run
- **WHEN** no `state.json` exists in `FOGHORN_PERSISTENT_DIR`
- **THEN** the container creates the file with the current hash, timestamp, and URL, and outputs status `pass` with `data.first_run: true`

#### Scenario: State file updated on no change
- **WHEN** the current content hash matches the stored hash
- **THEN** the timestamp in `state.json` is updated and the container outputs status `pass` with `data.changed: false`

#### Scenario: State file updated on change
- **WHEN** the current content hash differs from the stored hash
- **THEN** `state.json` is updated with the new hash and the container outputs status `fail` with `data.changed: true`

#### Scenario: Persistent directory not available
- **WHEN** `FOGHORN_PERSISTENT_DIR` is not set or the directory does not exist
- **THEN** the container outputs JSON with status `fail` and message indicating persistent memory is required

### Requirement: Change detection output
On change detection, the output SHALL include `data.previous_hash`, `data.current_hash`, and `data.changed: true`. On no change, `data.changed: false`. On first run, `data.first_run: true`.

#### Scenario: Change detected output
- **WHEN** the page content has changed since the last run
- **THEN** the output JSON includes `status: "fail"`, `message` describing the change, and `data` with `changed: true`, `previous_hash`, `current_hash`

#### Scenario: No change output
- **WHEN** the page content matches the stored state
- **THEN** the output JSON includes `status: "pass"`, `message` "No change detected", and `data` with `changed: false`, `current_hash`

#### Scenario: First run output
- **WHEN** the container runs for the first time with no stored state
- **THEN** the output JSON includes `status: "pass"`, `message` "Initial snapshot saved", and `data` with `first_run: true`, `current_hash`

### Requirement: Technical error handling
When the check cannot be performed due to a technical error (network failure, DNS resolution failure, connection refused, browser crash, SSL handshake failure when verification enabled, or timeout), the container SHALL output status `unknown`. This distinguishes "the check could not run" from "a change was detected" (`fail`). The `data` object SHALL include `error_type` (string identifying the error category) and `url`.

#### Scenario: Network / connectivity error
- **WHEN** the browser cannot reach the target URL (DNS error, connection refused, network unreachable)
- **THEN** the container outputs JSON with status `unknown`, message describing the error, and `data.error_type` set to `network`

#### Scenario: Browser navigation timeout
- **WHEN** the page does not load within the configured timeout
- **THEN** the container outputs JSON with status `unknown`, message indicating timeout, and `data.error_type` set to `timeout`

#### Scenario: SSL certificate failure
- **WHEN** `VERIFY_SSL` is `true` and the target has an invalid certificate
- **THEN** the container outputs JSON with status `unknown`, message describing the certificate error, and `data.error_type` set to `ssl`

#### Scenario: Browser crash or internal error
- **WHEN** the browser process crashes or Playwright raises an internal error
- **THEN** the container outputs JSON with status `unknown`, message describing the error, and `data.error_type` set to `browser`

### Requirement: Standard JSON output format
The container SHALL output JSON to stdout conforming to the Foghorn check contract: `{status, message, data, timestamp, duration_ms}`. The `data` object SHALL include `url`, `changed`, `current_hash`, `xpath` (when configured), `first_run` (on first run), and `previous_hash` (on change).

#### Scenario: Valid JSON output on change
- **WHEN** a change is detected
- **THEN** stdout is valid JSON with `status: "fail"`, all required fields, and change details in `data`

#### Scenario: Valid JSON output on no change
- **WHEN** no change is detected
- **THEN** stdout is valid JSON with `status: "pass"` and all required fields

### Requirement: SSL verification control
The container SHALL support `VERIFY_SSL` (default: `true`). When `false`, TLS certificate errors SHALL be ignored.

#### Scenario: SSL verification disabled
- **WHEN** `VERIFY_SSL` is `false`
- **THEN** the browser navigates to HTTPS pages with invalid certificates without error

#### Scenario: SSL verification enabled (default)
- **WHEN** `VERIFY_SSL` is not set or is `true`
- **THEN** the browser rejects pages with invalid certificates and outputs status `unknown` (see Technical error handling requirement)

### Requirement: Timeout configuration
The container SHALL respect `FOGHORN_TIMEOUT` (if set) as the maximum execution time. The default timeout SHALL be 30 seconds.

#### Scenario: FOGHORN_TIMEOUT respected
- **WHEN** `FOGHORN_TIMEOUT` is set to `15s`
- **THEN** the container enforces a 15-second overall timeout

### Requirement: Container directory structure
The container SHALL follow the standard Foghorn check container layout: `containers/web-change-check/` containing `Dockerfile`, `check.py`, `README.md`, and `VERSION`.

#### Scenario: Container files exist
- **WHEN** the web-change-check container is built
- **THEN** `containers/web-change-check/Dockerfile`, `containers/web-change-check/check.py`, `containers/web-change-check/README.md`, and `containers/web-change-check/VERSION` exist
