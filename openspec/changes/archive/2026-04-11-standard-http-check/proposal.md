## Why

The http-check container (`containers/http-check/`) is a minimal implementation that only supports basic URL fetching and single status code comparison. The legacy spec (`specs/standard-http-check.md`) defines a much richer feature set including status ranges, content regex, response time thresholds, custom methods/headers, and SSL verification control. Migrating to OpenSpec formalizes these requirements and creates a clear task list for bringing the implementation up to the full spec.

## What Changes

- Formalize the http-check capability as an OpenSpec change with full artifacts
- Document the gap between current implementation and the full spec requirements
- Bring the container implementation up to the spec: status ranges, content regex, response time thresholds, custom methods, headers, redirect control, and SSL toggle

## Capabilities

### New Capabilities
- `http-check`: HTTP/HTTPS endpoint verification via Docker container — validates status codes (single, range, list), response time thresholds, content regex, SSL verification, custom methods/headers, and redirect following

### Modified Capabilities
<!-- None -->

## Impact

- `containers/http-check/check.sh` — significant rewrite to support full feature set
- `containers/http-check/Dockerfile` — may need additional dependencies
- `openspec/specs/` — new spec directory for http-check
- `specs/standard-http-check.md` — legacy spec remains, superseded by OpenSpec artifacts
