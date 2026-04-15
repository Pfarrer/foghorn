# Web Change Check

Detects changes on web pages using a headless browser. Renders JavaScript, extracts content via XPath, and compares against previous runs.

Requires `persistent_memory: true` in check config.

## Image
`ghcr.io/pfarrer/foghorn-web-change-check:0.1.0`

## Env
- `URL` (required) - Target web page URL
- `XPATH` - XPath expression to scope content extraction (default: full body text)
- `WAIT_SECONDS` - Seconds to wait after page load (default: 0)
- `WAIT_FOR_SELECTOR` - CSS selector to wait for before extraction
- `VERIFY_SSL` - Verify TLS certificates (default: true)
- `FOGHORN_TIMEOUT` - Maximum execution time in seconds (default: 30)
- `FOGHORN_PERSISTENT_DIR` - Persistent storage path (required)
