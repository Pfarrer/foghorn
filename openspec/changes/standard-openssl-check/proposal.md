## Why

The openssl-check container exists as a working implementation in `containers/openssl-check/` with a legacy spec in `specs/standard-openssl-check.md`. Migrating to the OpenSpec workflow provides structured artifacts (proposal, design, specs, tasks) for ongoing maintenance, discoverability, and consistency with the project's evolving change management process.

## What Changes

- Formalize the existing openssl-check container as an OpenSpec change with full artifacts
- Port the legacy spec into the OpenSpec spec format under `openspec/specs/openssl-check/`
- No functional changes to the container or its behavior

## Capabilities

### New Capabilities
- `openssl-check`: TLS connection verification via Docker container — validates handshake, trust chain, hostname, certificate expiry, and minimum TLS version

### Modified Capabilities
<!-- None — this is a migration of an existing implemented feature -->

## Impact

- `openspec/specs/` — new spec directory and spec.md for openssl-check
- `specs/standard-openssl-check.md` — legacy spec remains but is superseded by OpenSpec artifacts
- No code changes to `containers/openssl-check/` or any Go packages
