# Basic Logging Implementation

## Category
functional

## Description
Implement structured logging system with configurable log levels via CLI arguments

## Usage Steps
1. Run with `-v` or `--verbose` flag for additional logs
2. Observe log output showing scheduler events, check executions, and errors

## Implementation Notes
- Use standard Go `log` package or structured logging library
- Log events should include:
  - Scheduler lifecycle (start/stop, tick processing)
  - Check execution (started, completed, queued, next run time)
  - Docker container operations (create, start, stop, remove)
  - Errors with context (check name, error details)
  - Timeout events
- Output format: `[LEVEL] message` (with optional timestamp in verbose mode)
- Log to stdout only

## Acceptance Criteria
- [ ] Debug level shows scheduler ticks and check queuing details
- [ ] Info level shows check start/completion and scheduler lifecycle events
- [ ] Warn level shows timeouts and warning conditions
- [ ] Error level shows failures and error conditions
- [ ] `-v/--verbose` flag adds timestamps
- [ ] Check execution logs include check name and status
- [ ] Error logs include relevant context (check name, error message)
- [ ] Timeout events are logged with check name and timeout duration
