## Why

Check execution history tracking is spec'd but not fully implemented. The state log (`state/log.go`) already persists basic records (check name, status, duration, timestamp) and the scheduler maintains in-memory history (last 10 entries). However, the spec requires start/end times, error details, and query-by-name-and-range capabilities that do not yet exist.

## What Changes

- Extend the state log record to include start time and error details
- Add query support for history by check name and time range
- Enforce retention policy based on record count or age
- Expose history via the status API

## Capabilities

### New Capabilities
- `check-execution-history`: Persistent check execution history with start/end times, durations, results, error details, and query support by check name and time range

### Modified Capabilities
<!-- None -->

## Impact

- `state/log.go` — extend `Record` struct with start time and error details
- `state/log.go` — add query methods (by name, by time range)
- `internal/statusapi/` — expose history in status API
- `config/types.go` — optional retention config
- `specs/check-execution-history.md` — legacy spec to be removed
