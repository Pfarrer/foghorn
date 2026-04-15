## 1. Container Skeleton

- [x] 1.1 Create `containers/web-change-check/` directory with `VERSION` (set to `0.1.0`), `README.md` (brief usage docs), and empty `check.py`
- [x] 1.2 Create `Dockerfile` using `mcr.microsoft.com/playwright/python` base image, installing dependencies and copying `check.py`

## 2. Core Check Script

- [x] 2.1 Implement env var parsing in `check.py`: `URL`, `XPATH`, `WAIT_SECONDS`, `WAIT_FOR_SELECTOR`, `VERIFY_SSL`, `FOGHORN_TIMEOUT`, `FOGHORN_PERSISTENT_DIR`
- [x] 2.2 Implement URL validation: reject missing/invalid URLs with `fail` status
- [x] 2.3 Implement headless browser launch and page navigation with timeout handling
- [x] 2.4 Implement content extraction: full body text when no XPath, XPath-matched nodes when `XPATH` is set
- [x] 2.5 Implement wait strategies: `WAIT_SECONDS` delay and `WAIT_FOR_SELECTOR` element wait
- [x] 2.6 Implement SHA-256 content hashing and state file read/write at `FOGHORN_PERSISTENT_DIR/state.json`

## 3. Change Detection Logic

- [x] 3.1 Implement first-run detection: no state file → save state, output `pass` with `first_run: true`
- [x] 3.2 Implement no-change detection: hash match → output `pass` with `changed: false`
- [x] 3.3 Implement change detection: hash mismatch → update state, output `fail` with `changed: true`, `previous_hash`, `current_hash`
- [x] 3.4 Implement missing `FOGHORN_PERSISTENT_DIR` guard: output `fail` with descriptive message

## 4. Output & Error Handling

- [x] 4.1 Ensure all code paths output valid Foghorn JSON (`{status, message, data, timestamp, duration_ms}`)
- [x] 4.2 Handle technical errors (navigation failure, timeout, browser crash, SSL failure) with `unknown` status and `data.error_type` field (`network`, `timeout`, `ssl`, `browser`)
- [x] 4.3 Implement SSL verification control via `VERIFY_SSL` env var

## 5. Verification

- [x] 5.1 Build the Docker image: `docker build -t foghorn-web-change-check containers/web-change-check/`
- [x] 5.2 Test first run (no state): verify `pass` output with `first_run: true` and state file created
- [x] 5.3 Test no-change run: run twice against same page, verify second run outputs `changed: false`
- [x] 5.4 Test change detection: modify content between runs, verify `fail` with `changed: true`
- [x] 5.5 Test XPath extraction with a known test page
- [x] 5.6 Test error cases: missing URL, missing `FOGHORN_PERSISTENT_DIR`, invalid XPath
- [x] 5.7 Test technical errors produce `unknown` status with `data.error_type` (simulate unreachable URL, timeout, SSL failure)
