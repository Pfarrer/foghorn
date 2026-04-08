## ADDED Requirements

### Requirement: Required input validation
The container SHALL validate that `HOST` and `PORT` are provided and that `PORT` is an integer between 1 and 65535. Missing or invalid inputs SHALL produce a `fail` status with a descriptive message.

#### Scenario: Missing HOST
- **WHEN** the `HOST` environment variable is not set
- **THEN** the container outputs JSON with status `fail` and message "HOST is required"

#### Scenario: Missing PORT
- **WHEN** the `PORT` environment variable is not set
- **THEN** the container outputs JSON with status `fail` and message "PORT is required"

#### Scenario: Invalid PORT
- **WHEN** `PORT` is set to a non-integer value or a value outside 1-65535
- **THEN** the container outputs JSON with status `fail` and message indicating PORT must be an integer between 1 and 65535

### Requirement: TLS handshake and trust chain verification
The container SHALL perform a TLS handshake using `openssl s_client` and verify the server certificate chain. If the CA bundle path is provided via `CA_BUNDLE_PATH`, it SHALL use that bundle; otherwise it SHALL use the system trust store.

#### Scenario: Successful trusted connection
- **WHEN** the target presents a valid certificate chain trusted by the system or custom CA bundle
- **THEN** the container outputs JSON with `trusted: true` and status `pass` or `warn` depending on expiry

#### Scenario: Untrusted certificate chain
- **WHEN** the target presents a certificate not trusted by the configured trust store
- **THEN** the container outputs JSON with status `fail` and message indicating the certificate is not trusted

#### Scenario: Custom CA bundle
- **WHEN** `CA_BUNDLE_PATH` is set to a valid file path
- **THEN** the container uses that file as the CA bundle for verification
- **WHEN** `CA_BUNDLE_PATH` is set to a non-existent path
- **THEN** the container outputs JSON with status `fail` indicating the path was not found

### Requirement: Hostname verification
The container SHALL verify the certificate hostname matches the target `HOST` by default. This behavior SHALL be controlled by the `VERIFY_HOSTNAME` environment variable (default: `true`). When enabled and the hostname does not match, the check SHALL fail.

#### Scenario: Hostname matches with verification enabled
- **WHEN** `VERIFY_HOSTNAME` is `true` and the certificate hostname matches `HOST`
- **THEN** `hostname_match` is `true` in the output data

#### Scenario: Hostname mismatch with verification enabled
- **WHEN** `VERIFY_HOSTNAME` is `true` and the certificate hostname does not match `HOST`
- **THEN** the container outputs JSON with status `fail` and message indicating hostname mismatch

#### Scenario: Hostname verification disabled
- **WHEN** `VERIFY_HOSTNAME` is `false`
- **THEN** hostname verification is skipped and the check does not fail due to hostname mismatch

### Requirement: Minimum TLS version selection
The container SHALL support selecting a minimum TLS version via the `MIN_TLS_VERSION` environment variable (default: `1.2`). Accepted values SHALL be `1.0`, `1.1`, `1.2`, and `1.3`. Invalid values SHALL produce a `fail` status.

#### Scenario: Minimum TLS version enforced
- **WHEN** `MIN_TLS_VERSION` is set to `1.3` and the server only supports TLS 1.2
- **THEN** the handshake fails and the container outputs status `fail`

#### Scenario: Invalid MIN_TLS_VERSION
- **WHEN** `MIN_TLS_VERSION` is set to an unrecognized value
- **THEN** the container outputs JSON with status `fail` indicating valid values

### Requirement: Certificate details reporting
The container SHALL report the certificate subject, issuer, validity dates (`not_before`, `not_after`), days remaining until expiry, TLS version negotiated, and cipher suite used.

#### Scenario: Successful certificate parse
- **WHEN** the TLS handshake succeeds and a peer certificate is available
- **THEN** the output `data` object includes `subject`, `issuer`, `not_before`, `not_after`, `days_remaining`, `tls_version`, and `cipher` fields

#### Scenario: Expired certificate
- **WHEN** the peer certificate's `not_after` date is in the past
- **THEN** `days_remaining` is `0` and status is `fail`

### Requirement: Expiry warning threshold
The container SHALL return status `warn` when the certificate is valid and trusted but expires within `WARNING_DAYS` days (default: 30).

#### Scenario: Certificate expiring within threshold
- **WHEN** the certificate is valid, trusted, and `days_remaining` is less than or equal to `WARNING_DAYS`
- **THEN** the container outputs status `warn` with message indicating days remaining

#### Scenario: Certificate expiring after threshold
- **WHEN** the certificate is valid, trusted, and `days_remaining` is greater than `WARNING_DAYS`
- **THEN** the container outputs status `pass`

### Requirement: Timeout handling
The container SHALL respect `TIMEOUT_SECONDS` (default: 10) and `FOGHORN_TIMEOUT` for connection timeouts. If `FOGHORN_TIMEOUT` is set and is shorter than `TIMEOUT_SECONDS`, it SHALL take precedence. Timeouts SHALL produce a `fail` status with a descriptive message.

#### Scenario: Connection timeout
- **WHEN** the TLS connection does not complete within the configured timeout
- **THEN** the container outputs JSON with status `fail` and message indicating the timeout

#### Scenario: FOGHORN_TIMEOUT overrides
- **WHEN** `FOGHORN_TIMEOUT` is set to a value shorter than `TIMEOUT_SECONDS`
- **THEN** the effective timeout is `FOGHORN_TIMEOUT`

### Requirement: Standard JSON output format
The container SHALL output JSON to stdout conforming to the Foghorn check contract: `{status, message, data, timestamp, duration_ms}`. The `data` object SHALL include `host`, `port`, `sni`, `tls_version`, `cipher`, `subject`, `issuer`, `not_before`, `not_after`, `days_remaining`, `trusted`, and `hostname_match`.

#### Scenario: Valid JSON output on success
- **WHEN** the check completes successfully
- **THEN** stdout is valid JSON with all required fields populated

#### Scenario: Valid JSON output on failure
- **WHEN** the check encounters an error
- **THEN** stdout is valid JSON with status `fail` or `unknown` and a descriptive message

### Requirement: SNI support
The container SHALL send the Server Name Indication (SNI) extension. The SNI value SHALL default to `HOST` but can be overridden via the `SNI` environment variable.

#### Scenario: Default SNI
- **WHEN** `SNI` is not set
- **THEN** the container uses the value of `HOST` as the SNI server name

#### Scenario: Custom SNI
- **WHEN** `SNI` is set to a different hostname
- **THEN** the container uses that value as the SNI server name
