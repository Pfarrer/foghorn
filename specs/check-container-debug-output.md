# Check Container Debug Output Modes

## Category
functional

## Description
Add configurable debug log output for check containers to simplify troubleshooting while keeping secret values protected.

## Usage Steps
1. Set debug output mode globally or per check (`off`, `on_failure`, `always`).
2. Run Foghorn with debug logging.
3. Inspect redacted container stdout/stderr in daemon logs.

## Implementation Notes
- Extend config with `debug_output` mode (global default and per-check override).
- Capture container stdout and stderr independently of result parsing.
- Log truncated output tails with a configurable size limit.
- Apply redaction before logging (`secret` values, secret-file contents, auth headers, password/token patterns).
- Keep current behavior as default (`off`) to avoid noisy logs.
- Add tests for:
  - mode behavior (`off`, `on_failure`, `always`)
  - truncation behavior
  - no secret leaks in logged debug output

## Acceptance Criteria
- [x] Config supports `debug_output: off|on_failure|always` globally and per check.
- [x] `on_failure` logs redacted container output only when check run fails.
- [x] `always` logs redacted container output for every check run.
- [x] Logged output is truncated to configured maximum size.
- [x] Tests verify mode behavior and secret redaction.

Passes: true
