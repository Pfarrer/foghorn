# HTTP Check

Runs an HTTP status check and returns Foghorn JSON output.

## Image
`ghcr.io/pfarrer/foghorn-http-check:1.0.2`

## Changelog
- 1.0.2 Add `ca-certificates` to support TLS certificate verification
- 1.0.1 Fix JSON output when curl request fails
- 1.0.0 Initial release

## Env
- `CHECK_URL`
- `EXPECTED_STATUS` (default `200`)
- `REQUEST_TIMEOUT` (default `10s`)
