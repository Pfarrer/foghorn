# OpenSSL TLS Check

Verifies TLS connectivity, trust chain, hostname, and certificate expiry.

## Image
`ghcr.io/pfarrer/foghorn-openssl-check:1.0.0`

## Changelog
- 1.0.0 Initial release

## Env
- `HOST` (required)
- `PORT` (required)
- `SNI` (default `HOST`)
- `MIN_TLS_VERSION` (default `1.2`)
- `CA_BUNDLE_PATH` (optional)
- `VERIFY_HOSTNAME` (default `true`)
- `WARNING_DAYS` (default `30`)
- `TIMEOUT_SECONDS` (default `10`)
- `FOGHORN_TIMEOUT` (optional; caps timeout)
