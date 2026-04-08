## Context

The openssl-check container (`containers/openssl-check/`) is a fully implemented Docker container that validates TLS connections. It currently exists as a working implementation backed by a legacy spec at `specs/standard-openssl-check.md`. This change formalizes it into the OpenSpec workflow without modifying any code.

The container uses a shell script (`check.sh`) on Debian Bookworm with the `openssl` CLI. It follows the standard Foghorn check contract: env vars in, JSON to stdout.

## Goals / Non-Goals

**Goals:**
- Document the existing openssl-check implementation as formal OpenSpec artifacts
- Preserve all current behavior in structured, testable spec requirements
- Enable future changes to this capability through the OpenSpec workflow

**Non-Goals:**
- Modifying the container implementation or its behavior
- Rewriting in Go or changing the base image
- Adding new features beyond what exists today

## Decisions

**Shell-based implementation on Debian over Alpine/Go**: The current implementation uses `check.sh` on `debian:bookworm-slim` for broad `openssl` compatibility and `date` parsing. This is retained as-is.

**Environment variable interface**: The container accepts `HOST`, `PORT`, `SNI`, `MIN_TLS_VERSION`, `CA_BUNDLE_PATH`, `VERIFY_HOSTNAME`, `WARNING_DAYS`, `TIMEOUT_SECONDS`, and `FOGHORN_TIMEOUT`. This interface is documented as the contract.

**JSON output format**: Standard Foghorn check output with `status`, `message`, `data`, `timestamp`, `duration_ms`. The `data` object includes TLS-specific fields: `host`, `port`, `sni`, `tls_version`, `cipher`, `subject`, `issuer`, `not_before`, `not_after`, `days_remaining`, `trusted`, `hostname_match`.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk; artifacts describe existing behavior
- [Legacy spec coexists] → Both `specs/standard-openssl-check.md` and OpenSpec specs will exist until legacy is retired; update `specs/STATUS.md` to reference OpenSpec location
