## 1. Extend State Log Records

- [x] 1.1 Add `started_at` (time.Time) and `error_details` (string) fields to `state.Record` struct
- [x] 1.2 Ensure backward compatibility: old records without new fields parse with zero values
- [x] 1.3 Update `RecordResult` to accept start time and error details

## 2. Query Methods

- [x] 2.1 Add `Query(checkName string, start, end time.Time) ([]Record, error)` method to `StateLog`
- [x] 2.2 Support filtering by check name only (omit time range)
- [x] 2.3 Support filtering by time range only (omit check name)
- [x] 2.4 Support combined check name + time range filtering

## 3. Status API Integration

- [x] 3.1 Expose history query endpoint in `internal/statusapi`
- [x] 3.2 Support query parameters: check name, start time, end time

## 4. Testing

- [x] 4.1 Add tests for extended record persistence and loading
- [x] 4.2 Add tests for query by check name
- [x] 4.3 Add tests for query by time range
- [x] 4.4 Add tests for combined filters
- [x] 4.5 Add tests for retention pruning with extended records
- [x] 4.6 Verify code compiles without warnings and all tests pass

## 5. Documentation

- [x] 5.1 Delete legacy spec file `specs/check-execution-history.md`
- [x] 5.2 Update `specs/STATUS.md` to remove the entry

## 6. Archive Change

- [x] 6.1 Archive the OpenSpec change once implementation is complete
