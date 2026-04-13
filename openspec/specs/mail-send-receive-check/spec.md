## ADDED Requirements

### Requirement: SMTP configuration validation
The container SHALL validate that `SMTP_HOST`, `SMTP_USERNAME`, `SMTP_PASSWORD` (or `SMTP_PASSWORD_FILE`), `MAIL_FROM`, and `MAIL_TO` are provided. Missing required values SHALL produce `unknown` status. `SMTP_PORT` SHALL default to 587 and validate as an integer 1-65535.

#### Scenario: Missing SMTP_HOST
- **WHEN** `SMTP_HOST` is not set
- **THEN** the container outputs status `unknown` listing the required SMTP fields

#### Scenario: Invalid SMTP_PORT
- **WHEN** `SMTP_PORT` is not a valid integer between 1 and 65535
- **THEN** the container outputs status `unknown` indicating invalid port

#### Scenario: Secret file fallback for SMTP_PASSWORD
- **WHEN** `SMTP_PASSWORD` is not set but `SMTP_PASSWORD_FILE` points to an existing file
- **THEN** the container reads the password from that file

### Requirement: SMTP TLS mode
The container SHALL support `SMTP_TLS_MODE` with values `starttls` (default), `tls`, and `none`. `starttls` SHALL use `--ssl-reqd`; `tls` SHALL use `smtps://`; `none` SHALL use plaintext SMTP. Invalid values SHALL produce `unknown`.

#### Scenario: StartTLS mode
- **WHEN** `SMTP_TLS_MODE` is `starttls`
- **THEN** curl connects via `smtp://` with `--ssl-reqd`

#### Scenario: Implicit TLS mode
- **WHEN** `SMTP_TLS_MODE` is `tls`
- **THEN** curl connects via `smtps://`

#### Scenario: No TLS mode
- **WHEN** `SMTP_TLS_MODE` is `none`
- **THEN** curl connects via `smtp://` without SSL flags

#### Scenario: Invalid TLS mode
- **WHEN** `SMTP_TLS_MODE` is set to an unrecognized value
- **THEN** the container outputs status `unknown`

### Requirement: Unique message generation
The container SHALL generate a unique correlation ID embedded in the email subject and body. The subject SHALL be `{SUBJECT_PREFIX} {correlation_id}`. The body SHALL include the correlation ID and send timestamp. A `Message-ID` header SHALL be set.

#### Scenario: Correlation ID format
- **WHEN** a message is generated
- **THEN** the subject contains a unique correlation ID matching the pattern `foghorn-<timestamp>-<random_hex>`

#### Scenario: Body contains tracking data
- **WHEN** the message body is constructed
- **THEN** it includes `correlation_id=` and `sent_at=` lines

### Requirement: SMTP message sending
The container SHALL send the generated message via curl's SMTP support. Send failure SHALL produce `fail` status. The remaining deadline time SHALL be used as the curl timeout.

#### Scenario: Successful send
- **WHEN** SMTP credentials and host are valid
- **THEN** the message is sent and the container proceeds to polling

#### Scenario: Send failure
- **WHEN** SMTP authentication fails or the connection is refused
- **THEN** the container outputs status `fail` with message "failed to send mail via SMTP"

### Requirement: IMAP receive configuration validation
The container SHALL validate that `RECEIVE_HOST`, `RECEIVE_USERNAME`, and `RECEIVE_PASSWORD` (or `RECEIVE_PASSWORD_FILE`) are provided. `RECEIVE_PORT` SHALL default to 993. `RECEIVE_TLS` SHALL default to `true`.

#### Scenario: Missing receive config
- **WHEN** `RECEIVE_HOST`, `RECEIVE_USERNAME`, or `RECEIVE_PASSWORD` is not set
- **THEN** the container outputs status `unknown` listing required receive fields

#### Scenario: Secret file fallback for RECEIVE_PASSWORD
- **WHEN** `RECEIVE_PASSWORD` is not set but `RECEIVE_PASSWORD_FILE` points to an existing file
- **THEN** the container reads the password from that file

### Requirement: IMAP polling with deadline
The container SHALL poll the IMAP mailbox searching for the correlation ID in the subject. Polling SHALL continue at `POLL_INTERVAL_SECONDS` (default: 5) intervals until `DEADLINE_SECONDS` expires. If `FOGHORN_TIMEOUT` is set and shorter, it overrides the deadline. Expired deadline without match SHALL produce `fail`.

#### Scenario: Message found before deadline
- **WHEN** the correlation ID is found in IMAP search results within the deadline
- **THEN** the container proceeds to process the matched message

#### Scenario: Deadline expires without match
- **WHEN** the deadline expires and no matching message is found
- **THEN** the container outputs status `fail` indicating the message did not arrive in time

#### Scenario: IMAP search failure
- **WHEN** the IMAP search request fails
- **THEN** the container outputs status `fail` indicating failure to query the mailbox

### Requirement: Stale message filtering
The container SHALL ignore matched messages with an INTERNALDATE older than the check start time, preventing false positives from previous runs.

#### Scenario: Old message ignored
- **WHEN** a message matches the correlation ID but has INTERNALDATE before the check started
- **THEN** the message is ignored and polling continues

#### Scenario: Fresh message accepted
- **WHEN** a message matches the correlation ID and has INTERNALDATE at or after the check start time
- **THEN** the message is accepted as the delivery result

### Requirement: Delivery time calculation
The container SHALL calculate delivery time as the difference between the receive timestamp (INTERNALDATE or current time) and the send timestamp. Negative values SHALL be clamped to 0.

#### Scenario: Normal delivery
- **WHEN** a message is sent at T0 and received at T0+18s
- **THEN** `delivery_seconds` is 18

#### Scenario: Delivery time clamped
- **WHEN** the calculated delivery time is negative due to clock skew
- **THEN** `delivery_seconds` is 0

### Requirement: Warning threshold
The container SHALL support optional `WARNING_THRESHOLD_SECONDS`. When set and delivery time exceeds this value, status SHALL be `warn` instead of `pass`.

#### Scenario: Delivery exceeds warning threshold
- **WHEN** `WARNING_THRESHOLD_SECONDS` is 15 and delivery takes 18s
- **THEN** status is `warn` with message indicating the threshold was exceeded

#### Scenario: Delivery within warning threshold
- **WHEN** `WARNING_THRESHOLD_SECONDS` is 30 and delivery takes 18s
- **THEN** status is `pass`

#### Scenario: No warning threshold set
- **WHEN** `WARNING_THRESHOLD_SECONDS` is not set
- **THEN** status is `pass` for any delivery time within the deadline

### Requirement: Delete after match
The container SHALL support `DELETE_AFTER_MATCH` (default: `false`). When `true`, the matched message SHALL be marked as deleted and expunged from the mailbox.

#### Scenario: Delete enabled
- **WHEN** `DELETE_AFTER_MATCH` is `true` and a message is matched
- **THEN** the container issues UID STORE +FLAGS (\Deleted) and EXPUNGE

#### Scenario: Delete disabled
- **WHEN** `DELETE_AFTER_MATCH` is `false` (default)
- **THEN** the matched message is left in the mailbox

### Requirement: Standard JSON output format
The container SHALL output JSON to stdout conforming to the Foghorn check contract: `{status, message, data, timestamp, duration_ms}`. The `data` object SHALL include `mail_from`, `mail_to`, `smtp_host`, `receive_host`, `receive_protocol`, `message_id`, `correlation_id`, `send_time`, `receive_time`, `delivery_seconds`, and `deadline_seconds`.

#### Scenario: Valid JSON on success
- **WHEN** the check completes successfully
- **THEN** stdout is valid JSON with all required fields populated

#### Scenario: Valid JSON on failure
- **WHEN** the check fails at any step
- **THEN** stdout is valid JSON with status `fail` or `unknown` and empty `receive_time`

### Requirement: Timeout handling
The container SHALL respect `DEADLINE_SECONDS` as the overall check deadline. `FOGHORN_TIMEOUT` SHALL override the deadline if shorter. All curl operations SHALL use remaining deadline time as their max timeout.

#### Scenario: FOGHORN_TIMEOUT overrides deadline
- **WHEN** `FOGHORN_TIMEOUT` is set to 30s and `DEADLINE_SECONDS` is 60s
- **THEN** the effective deadline is 30s

#### Scenario: Curl timeouts bounded by remaining deadline
- **WHEN** the remaining deadline is 20s
- **THEN** curl operations use at most 20s as their max timeout
