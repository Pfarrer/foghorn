# Mail Send/Receive Check

Sends a probe mail through SMTP and verifies delivery by polling IMAP.

## Image
`ghcr.io/pfarrer/foghorn-mail-send-receive-check:1.0.0`

## Changelog
- 1.0.0 Initial release

## Env
- `SMTP_HOST` (required)
- `SMTP_PORT` (default `587`)
- `SMTP_USERNAME` (required)
- `SMTP_PASSWORD` (required)
- `SMTP_PASSWORD_FILE` (optional; file path alternative to `SMTP_PASSWORD`)
- `SMTP_TLS_MODE` (default `starttls`; `starttls|tls|none`)
- `MAIL_FROM` (required)
- `MAIL_TO` (required)
- `SUBJECT_PREFIX` (default `FOGHORN-CHECK`)
- `BODY_TEMPLATE` (optional)
- `RECEIVE_HOST` (required)
- `RECEIVE_PORT` (default `993`)
- `RECEIVE_USERNAME` (required)
- `RECEIVE_PASSWORD` (required)
- `RECEIVE_PASSWORD_FILE` (optional; file path alternative to `RECEIVE_PASSWORD`)
- `RECEIVE_TLS` (default `true`)
- `RECEIVE_MAILBOX` (default `INBOX`)
- `POLL_INTERVAL_SECONDS` (default `5`)
- `WARNING_THRESHOLD_SECONDS` (optional)
- `DEADLINE_SECONDS` (required)
- `DELETE_AFTER_MATCH` (default `false`)
- `FOGHORN_TIMEOUT` (optional; caps deadline)
