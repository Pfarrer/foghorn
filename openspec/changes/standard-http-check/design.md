## Context

The http-check container (`containers/http-check/`) currently performs a basic curl-based HTTP GET, compares a single status code, and outputs JSON. The legacy spec (`specs/standard-http-check.md`) defines a significantly richer feature set that is not yet implemented: status code ranges/lists, response time thresholds (warn/critical), content regex validation, custom HTTP methods, request headers, request body, redirect control, and detailed response data.

The current implementation uses `check.sh` on `debian:bookworm-slim` with `curl`. This is the right foundation — curl supports all the required features natively via flags.

## Goals / Non-Goals

**Goals:**
- Implement the full http-check spec: status ranges/lists, response time thresholds, content regex, custom methods, headers, body, redirect control, SSL toggle
- Output detailed response data: url, method, status_code, status_text, response_time_ms, response_size_bytes, content_match
- Follow the Foghorn check contract (env vars in, JSON to stdout)

**Non-Goals:**
- Changing the base image or language (keep shell + curl on Debian)
- Adding authentication handling beyond what headers provide
- Supporting HTTP/2-specific features beyond what curl handles

## Decisions

**Shell + curl over Go implementation**: curl already supports all required features (methods, headers, body, redirects, SSL, status code capture, timing, regex output matching). Shell keeps the container simple and consistent with other check containers (openssl-check, disk-check).

**Status code matching logic**: Parse `EXPECTED_STATUS` to support three formats: single code (`200`), range (`200-299`), comma-separated list (`200,201,204`). Range parsing done in shell arithmetic.

**Response time thresholds**: Use curl's `-w` write-out format to capture `time_total` and convert to ms. Compare against `WARNING_THRESHOLD_MS` (default 1000) and `CRITICAL_THRESHOLD_MS` (default 5000) after status code passes.

**Content regex**: Use curl to save response body to a temp file, then match with grep. Report `content_match: true/false` in data.

**Redirect control**: `FOLLOW_REDIRECTS` (default true) maps to curl's `-L` flag. When false, omit `-L`.

**SSL verification**: `VERIFY_SSL` (default true) maps to curl behavior. When false, add `-k` (insecure).

**URL env var naming**: The current implementation uses `CHECK_URL` but the spec uses `URL`. Accept both with `URL` taking precedence for backward compatibility.

## Risks / Trade-offs

- [Shell JSON escaping] → Use careful quoting and sed escaping; test with special characters in URLs/headers
- [curl version compatibility] → Debian Bookworm's curl is recent enough for all needed write-out variables
- [Backward compatibility] → Current users may rely on `CHECK_URL` env var; continue accepting it as fallback
