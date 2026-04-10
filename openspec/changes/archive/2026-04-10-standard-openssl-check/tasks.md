## 1. Verify Existing Implementation Against Specs

- [x] 1.1 Verify the openssl-check container produces valid JSON output matching the spec format (status, message, data, timestamp, duration_ms)
- [x] 1.2 Verify input validation: HOST required, PORT required and in range 1-65535, invalid values produce `fail`
- [x] 1.3 Verify TLS handshake and trust chain: trusted cert → pass/warn, untrusted → fail, custom CA bundle path support
- [x] 1.4 Verify hostname verification: default enabled, can be disabled via VERIFY_HOSTNAME=false, mismatch produces fail
- [x] 1.5 Verify MIN_TLS_VERSION enforcement: accepts 1.0/1.1/1.2/1.3, rejects invalid values, enforces on handshake
- [x] 1.6 Verify certificate details: subject, issuer, not_before, not_after, days_remaining, tls_version, cipher reported correctly
- [x] 1.7 Verify expiry warning: warn when days_remaining <= WARNING_DAYS, pass when above threshold
- [x] 1.8 Verify timeout handling: TIMEOUT_SECONDS default 10, FOGHORN_TIMEOUT override, timeout produces fail
- [x] 1.9 Verify SNI support: defaults to HOST, override via SNI env var

## 2. Update Documentation

- [x] 2.1 Update specs/STATUS.md to reference OpenSpec artifacts for standard-openssl-check
- [x] 2.2 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
