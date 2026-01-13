# Foghorn

A service monitoring tool that executes arbitrary Docker containers as health checks.

Checks are Docker containers that run on a predefined schedule, perform custom actions, and return results for evaluation. Unlike traditional monitoring that only supports ping or HTTP checks, Foghorn allows any containerized check.

## Features

- Schedule-based execution of Docker containers
- Arbitrary check logic through containerization
- Result evaluation and status reporting
- Cron and interval-based scheduling
- YAML configuration format

## Usage

Foghorn can be run with a configuration file:

```bash
./foghorn example.yaml
```

The scheduler will load the configuration and execute checks based on their cron schedules.

## Docker Check Interface

Foghorn executes Docker containers as health checks and communicates with them through a well-defined interface.

### Environment Variables

Check containers receive the following environment variables from Foghorn:

- `FOGHORN_CHECK_NAME`: Name of the check
- `FOGHORN_CHECK_CONFIG`: JSON string with check-specific configuration (from metadata)
- `FOGHORN_ENDPOINT`: Target endpoint to check (if applicable)
- `FOGHORN_SECRETS`: JSON string with secrets (for environment variables starting with `SECRET_`)
- `FOGHORN_TIMEOUT`: Timeout duration for the check

All other custom environment variables defined in the check configuration are also passed to the container.

### Output Format

Check containers should output JSON results to stdout with the following structure:

```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Human-readable result description",
  "data": {
    "custom_key": "custom_value"
  },
  "timestamp": "2025-01-13T12:00:00Z",
  "duration_ms": 150
}
```

Fields:
- `status` (required): Check status - "pass", "fail", "warn", or "unknown"
- `message` (required): Human-readable description of the result
- `data` (optional): Structured data/metrics from the check
- `timestamp` (required): ISO 8601 timestamp of when the check completed
- `duration_ms` (required): Check execution duration in milliseconds

### Exit Codes

- `0`: Check completed successfully (use status field in JSON for pass/fail)
- `non-zero`: Check encountered an error during execution

### Output Location

By default, Foghorn reads JSON output from container stdout. Alternatively, containers can write results to `/output/result.json` inside the container, and Foghorn will read that file if stdout parsing fails.

### Example Check Container

See `examples/docker-check/` for a complete example of a Docker container that implements the Foghorn check interface.

## Configuration

Foghorn uses YAML configuration files to define checks. See `example.yaml` for a comprehensive example of all available configuration options.

Configuration includes:
- Check definitions with container images
- Schedules (cron expressions or intervals)
- Evaluation rules for results
- Metadata and tags
- Environment variables and timeouts

## Scheduler

The scheduler component manages check execution based on cron expressions:
- Parses standard cron expressions (minute, hour, day, month, day of week)
- Calculates next execution time for each check
- Triggers check execution when scheduled time is reached
- Supports time zones for accurate scheduling
- Only executes enabled checks
