## Context

The mail-send-receive-check container (`containers/mail-send-receive-check/`) is a fully implemented Docker container that performs end-to-end email delivery verification. It sends a uniquely tagged message via SMTP, then polls an IMAP mailbox until the message arrives or a deadline expires. Runs on Debian Bookworm with curl for both SMTP and IMAP.

## Goals / Non-Goals

**Goals:**
- Document the existing implementation as formal OpenSpec artifacts
- Preserve all current behavior in structured, testable spec requirements
- Enable future changes through the OpenSpec workflow

**Non-Goals:**
- Modifying the container implementation
- Adding POP3 support or other receive protocols
- Adding multi-recipient or attachment support

## Decisions

**curl for SMTP and IMAP**: The implementation uses curl's built-in SMTP send (`--mail-from`, `--mail-rcpt`, `--upload-file`) and IMAP commands (`UID SEARCH`, `UID FETCH`, `UID STORE`, `EXPUNGE`). This avoids needing language runtimes or mail libraries.

**Correlation ID matching**: Each probe generates a unique correlation ID (`foghorn-<timestamp>-<random_hex>`) embedded in the subject and body. The receiver searches by subject header match, avoiding stale message confusion.

**Stale message filtering**: Messages found via IMAP with an INTERNALDATE older than the check start time are ignored, preventing false positives from previous runs.

**Secret file support**: `SMTP_PASSWORD` and `RECEIVE_PASSWORD` support `_FILE` suffixed variants (`SMTP_PASSWORD_FILE`, `RECEIVE_PASSWORD_FILE`) for Foghorn's secret injection pattern.

**TLS mode support**: SMTP supports three modes: `starttls` (default, uses `--ssl-reqd`), `tls` (implicit TLS via `smtps://`), and `none`. IMAP TLS controlled by `RECEIVE_TLS` boolean, switching between `imaps://` and `imap://`.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk; artifacts describe existing behavior
- [IMAP-only receive] → POP3 not supported; acceptable per current spec scope
