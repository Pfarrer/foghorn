# Standard Mail Send/Receive Check Container

## Category
integration

## Description
Define a reusable Docker container that sends a test mail through server A and verifies receipt on server B within a deadline.

## Usage Steps
1. User pulls the standard mail check image: `foghorn/mail-send-receive-check:1.0.0`.
2. Configure SMTP connection for server A and mailbox access for server B.
3. Foghorn runs the container on schedule.
4. Container sends a uniquely tagged mail, polls server B, and outputs results.

## Implementation Notes
Container should accept these environment variables:
- `SMTP_HOST` (required): Mail server A host.
- `SMTP_PORT` (default: 587)
- `SMTP_USERNAME` (required)
- `SMTP_PASSWORD` (required)
- `SMTP_TLS_MODE` (default: `starttls`; allowed: `starttls`, `tls`, `none`)
- `MAIL_FROM` (required): Sender address.
- `MAIL_TO` (required): Target mailbox on server B.
- `SUBJECT_PREFIX` (default: `FOGHORN-CHECK`)
- `BODY_TEMPLATE` (optional): Body text; implementation appends unique check id.
- `RECEIVE_HOST` (required): Mail server B host.
- `RECEIVE_PORT` (default: 993)
- `RECEIVE_USERNAME` (required)
- `RECEIVE_PASSWORD` (required)
- `RECEIVE_TLS` (default: `true`)
- `RECEIVE_MAILBOX` (default: `INBOX`)
- `POLL_INTERVAL_SECONDS` (default: 5)
- `WARNING_THRESHOLD_SECONDS` (optional): Return `warn` when delivery time exceeds this value.
- `DEADLINE_SECONDS` (required): Max wait time for arrival.
- `DELETE_AFTER_MATCH` (default: `false`): Delete matched message if supported.

Output JSON format:
```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Mail delivered in 18s",
  "data": {
    "mail_from": "probe@a.example.com",
    "mail_to": "probe@b.example.com",
    "smtp_host": "smtp.a.example.com",
    "receive_host": "imap.b.example.com",
    "receive_protocol": "imap",
    "message_id": "<...>",
    "correlation_id": "foghorn-2026-01-18T12:00:00Z-abc123",
    "send_time": "2026-01-18T12:00:00Z",
    "receive_time": "2026-01-18T12:00:18Z",
    "delivery_seconds": 18,
    "deadline_seconds": 60
  },
  "timestamp": "2026-01-18T12:00:18Z",
  "duration_ms": 18000
}
```

Status determination:
- `pass`: Message is sent and found on server B within `DEADLINE_SECONDS`.
- `fail`: Send fails, auth fails, connect fails, or message not found before deadline.
- `warn`: Message arrives but delivery time exceeds `WARNING_THRESHOLD_SECONDS` when set.
- `unknown`: Invalid config or unexpected runtime error.

Matching rules:
- Generated subject must include a unique correlation id.
- Receiver search must match correlation id (subject and/or body) to avoid stale message matches.
- Receiver should ignore messages older than check start time.

## Acceptance Criteria
- [ ] Validates required SMTP and receive configuration.
- [ ] Sends one uniquely tagged message through server A.
- [ ] Polls server B until matched message is found or deadline expires.
- [ ] Passes only when matched message arrives within `DEADLINE_SECONDS`.
- [ ] Fails when deadline expires without a match.
- [ ] Supports IMAP receive mode.
- [ ] Returns structured JSON with send/receive timestamps and delivery duration.
- [ ] Respects `FOGHORN_TIMEOUT` if provided.

## Example Configuration
```yaml
name: "mail-send-receive"
image: "foghorn/mail-send-receive-check:1.0.0"
schedule:
  interval: "5m"
env:
  SMTP_HOST: "smtp.server-a.example.com"
  SMTP_PORT: "587"
  SMTP_USERNAME: "probe-a"
  SMTP_PASSWORD: "${SMTP_PASSWORD}"
  SMTP_TLS_MODE: "starttls"
  MAIL_FROM: "probe@server-a.example.com"
  MAIL_TO: "probe@server-b.example.com"
  RECEIVE_HOST: "imap.server-b.example.com"
  RECEIVE_PORT: "993"
  RECEIVE_USERNAME: "probe-b"
  RECEIVE_PASSWORD: "${IMAP_PASSWORD}"
  RECEIVE_TLS: "true"
  POLL_INTERVAL_SECONDS: "5"
  DEADLINE_SECONDS: "60"
timeout: "90s"
```
