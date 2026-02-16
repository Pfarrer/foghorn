# Standard OpenSSL TLS Check Container

## Category
security

## Description
Define a reusable Docker container that verifies a trusted TLS connection to a host and port

## Usage Steps
1. User pulls the standard OpenSSL check image: `foghorn/openssl-check:1.0.0`
2. Configure check with target host and port
3. Foghorn runs the container on schedule
4. Container validates TLS handshake and trust chain, then outputs results

## Implementation Notes
Container should accept these environment variables:
- `HOST` (required): Target hostname
- `PORT` (required): Target port
- `SNI`: Server name for SNI (default: HOST)
- `MIN_TLS_VERSION`: Minimum TLS version (default: 1.2)
- `CA_BUNDLE_PATH`: Optional CA bundle file path; if unset, use system trust store
- `VERIFY_HOSTNAME`: Whether to verify certificate hostname (default: true)
- `WARNING_DAYS`: Days before expiry to return `warn` (default: 30)
- `TIMEOUT_SECONDS`: Connection timeout (default: 10)

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "TLS ok for example.com:443 (expires in 120d)",
  "data": {
    "host": "example.com",
    "port": 443,
    "sni": "example.com",
    "tls_version": "TLSv1.3",
    "cipher": "TLS_AES_256_GCM_SHA384",
    "subject": "CN=example.com",
    "issuer": "CN=Example CA",
    "not_before": "2025-01-01T00:00:00Z",
    "not_after": "2026-01-01T00:00:00Z",
    "days_remaining": 120,
    "trusted": true,
    "hostname_match": true
  },
  "timestamp": "2025-01-18T12:00:00Z",
  "duration_ms": 200
}
```

Status determination:
- `pass`: TLS handshake succeeds, chain trusted, hostname matches, days_remaining > WARNING_DAYS
- `warn`: TLS handshake succeeds, chain trusted, hostname matches, days_remaining <= WARNING_DAYS
- `fail`: Handshake failure, untrusted chain, hostname mismatch, expired cert, or invalid config
- `unknown`: Execution error or unexpected output

Use Alpine with `openssl` or Go-based TLS to perform the check

## Acceptance Criteria
- [ ] Validates `HOST` and `PORT` are provided
- [ ] Performs TLS handshake and verifies trust chain
- [ ] Verifies hostname by default with option to disable
- [ ] Supports minimum TLS version selection
- [ ] Reports certificate subject, issuer, and validity dates
- [ ] Calculates and reports days remaining to expiry
- [ ] Returns `warn` when certificate is nearing expiration
- [ ] Respects `FOGHORN_TIMEOUT` if provided
- [ ] Handles timeouts and connection errors cleanly

## Example Configuration
```yaml
name: "tls-check"
image: "foghorn/openssl-check:1.0.0"
schedule:
  interval: "5m"
env:
  HOST: "example.com"
  PORT: "443"
  WARNING_DAYS: "21"
timeout: "15s"
```

```yaml
name: "tls-check-custom-ca"
image: "foghorn/openssl-check:1.0.0"
schedule:
  cron: "0 */6 * * *"
env:
  HOST: "internal.example.com"
  PORT: "8443"
  CA_BUNDLE_PATH: "/etc/ssl/custom-ca.pem"
  MIN_TLS_VERSION: "1.3"
```

## Passes
false
