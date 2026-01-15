# Concurrent Execution and Timeout Handling

## Category
functional

## Description
Manage concurrent check executions with proper limits and implement timeouts for stuck or long-running checks

## Usage Steps
1. Configure maximum concurrent check executions
2. Define check-specific timeouts in configuration
3. Start the Foghorn service
4. Scheduler manages concurrent execution automatically
5. Stuck checks are terminated after timeout

## Implementation Notes
- Implement concurrency limits for check executions
- Queue checks when concurrency limit is reached
- Support configurable global concurrency limit
- Support per-check timeout configuration
- Implement timeout mechanism using context or signal
- Terminate Docker containers that exceed timeout
- Log timeout events appropriately

## Acceptance Criteria
- [x] Multiple checks can run concurrently up to configured limit
- [x] Checks are queued when concurrency limit is reached
- [x] Long-running checks have timeout and are terminated
- [x] Check execution is logged when timeout occurs
- [x] Timeout is configurable globally and per-check
- [x] No resource exhaustion occurs with many concurrent checks
