## 1. Verify SMTP Configuration and Sending

- [x] 1.1 Verify SMTP config validation: missing SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, MAIL_FROM, or MAIL_TO → unknown
- [x] 1.2 Verify SMTP_PORT validation: non-integer or out of range → unknown; default 587 when unset
- [x] 1.3 Verify SMTP_PASSWORD_FILE fallback: reads password from file when SMTP_PASSWORD is unset
- [x] 1.4 Verify SMTP_TLS_MODE: starttls uses --ssl-reqd, tls uses smtps://, none uses plaintext, invalid → unknown
- [x] 1.5 Verify unique correlation ID generation in subject and body
- [x] 1.6 Verify SMTP send failure produces fail status

## 2. Verify IMAP Receive and Polling

- [x] 2.1 Verify receive config validation: missing RECEIVE_HOST, RECEIVE_USERNAME, RECEIVE_PASSWORD → unknown
- [x] 2.2 Verify RECEIVE_PASSWORD_FILE fallback
- [x] 2.3 Verify IMAP polling searches by correlation ID in subject header
- [x] 2.4 Verify deadline expiry without match → fail
- [x] 2.5 Verify IMAP search failure → fail
- [x] 2.6 Verify stale message filtering: INTERNALDATE before check start → ignored
- [x] 2.7 Verify delivery_seconds calculated correctly, clamped to 0 on negative
- [x] 2.8 Verify RECEIVE_TLS: true → imaps://, false → imap://

## 3. Verify Thresholds and Post-Processing

- [x] 3.1 Verify WARNING_THRESHOLD_SECONDS: delivery exceeds threshold → warn; within → pass; unset → always pass
- [x] 3.2 Verify DELETE_AFTER_MATCH: true → UID STORE +FLAGS (\Deleted) + EXPUNGE; false → no deletion
- [x] 3.3 Verify FOGHORN_TIMEOUT overrides DEADLINE_SECONDS when shorter

## 4. Verify JSON Output

- [x] 4.1 Verify JSON output contains all required fields: status, message, data (mail_from, mail_to, smtp_host, receive_host, receive_protocol, message_id, correlation_id, send_time, receive_time, delivery_seconds, deadline_seconds), timestamp, duration_ms
- [x] 4.2 Verify JSON on failure has empty receive_time and valid structure

## 5. Documentation

- [x] 5.1 Delete legacy spec file `specs/standard-mail-send-receive-check.md`
- [x] 5.2 Update `specs/STATUS.md` to remove the standard-mail-send-receive-check entry
- [x] 5.3 Verify all Go tests pass and code compiles without warnings

## 6. Archive Change

- [x] 6.1 Archive the OpenSpec change once verification is complete
