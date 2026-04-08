## Why

The mail-send-receive-check container exists as a fully implemented Docker container in `containers/mail-send-receive-check/` with a legacy spec in `specs/standard-mail-send-receive-check.md`. Migrating to OpenSpec provides structured artifacts for ongoing maintenance and consistency with the project's change management process.

## What Changes

- Formalize the existing mail-send-receive-check container as an OpenSpec change with full artifacts
- Port the legacy spec into the OpenSpec spec format under `openspec/specs/mail-send-receive-check/`
- No functional changes to the container or its behavior

## Capabilities

### New Capabilities
- `mail-send-receive-check`: End-to-end mail delivery verification via Docker container — sends a uniquely tagged email through SMTP server A, polls IMAP server B for arrival, and reports delivery time with deadline enforcement

### Modified Capabilities
<!-- None — this is a migration of an existing implemented feature -->

## Impact

- `openspec/specs/` — new spec directory and spec.md for mail-send-receive-check
- `specs/standard-mail-send-receive-check.md` — legacy spec to be removed
- No code changes to `containers/mail-send-receive-check/` or any Go packages
