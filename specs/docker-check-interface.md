# Docker Check Interface Contract

## Category
integration

## Description
Define the interface contract between Foghorn and Docker check containers, including configuration inputs and expected output format

## Usage Steps
1. Design Docker container to accept environment variables for configuration
2. Container performs its check logic
3. Container outputs JSON result to stdout or specific file
4. Foghorn parses output and evaluates result

## Implementation Notes
- Define environment variable names for check configuration:
  - `FOGHORN_CHECK_NAME`: Name of the check
  - `FOGHORN_CHECK_CONFIG`: JSON string with check-specific configuration
  - `FOGHORN_ENDPOINT`: Target endpoint to check (if applicable)
  - `FOGHORN_SECRETS`: JSON string with secrets (API keys, tokens)
  - `FOGHORN_TIMEOUT`: Timeout duration for the check
- Define expected output format (JSON):
  - `status`: "pass", "fail", "warn", or "unknown"
  - `message`: Human-readable result description
  - `data`: Optional object with structured data/metrics
  - `timestamp`: ISO 8601 timestamp of the check
  - `duration_ms`: Check execution duration in milliseconds
- Support output to stdout (default) or file at `/output/result.json`
- Define error handling: non-zero exit code means check failure
- Document interface in README or specification file for check authors

## Acceptance Criteria
- [ ] Environment variable contract is documented
- [ ] Expected JSON output format is defined and validated
- [ ] Example Docker container implements the interface
- [ ] Foghorn correctly passes configuration via environment variables
- [ ] Foghorn parses stdout JSON output correctly
- [ ] Non-zero exit codes are treated as failures
- [ ] Missing or invalid JSON output is handled gracefully
- [ ] Interface documentation is clear for external check authors

## Passes
false
