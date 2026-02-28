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

Daemon run:

```bash
./foghorn-daemon -c example.yaml
```

TUI client run:

```bash
./foghorn-tui
```

### Run Daemon With Docker Compose

Start the daemon:

```bash
docker compose up -d
```

Tail logs:

```bash
docker compose logs -f foghorn
```

Stop the daemon:

```bash
docker compose down
```

The included `docker-compose.yml` publishes the status API on `7676` and starts the daemon with `--status-listen 0.0.0.0:7676`.

### Connect TUI To Compose Daemon

Run the TUI on the host and point it at the daemon status API:

```bash
./foghorn-tui --status-url http://127.0.0.1:7676
```

If the daemon is running on another host, replace the URL:

```bash
./foghorn-tui --status-url http://<daemon-host>:7676
```

### Command-Line Options

Foghorn supports the following command-line flags:

- `-c, --config <path>`: Path to the configuration file (required)
- `-l, --log-level <level>`: Set log level (debug, info, warn, error) (default: info)
- `-v, --verbose`: Enable verbose logging with timestamps and source file locations
- `-d, --dry-run`: Validate configuration only without running the scheduler
- `-i, --verify-image-availability`: Verify all Docker images in config are available locally
- `--status-listen <addr>`: Status API listen address (default: `127.0.0.1:7676`)
- `--state-log-file <path>`: Persist check results to a state log file
- `--secret-store-file <path>`: Path to encrypted secret store file
- `-h, --help`: Display help message and usage information

### Examples

Run with default settings:
```bash
./foghorn-daemon -c example.yaml
```

Run with debug logging:
```bash
./foghorn-daemon -c example.yaml -l debug
```

Run with verbose output:
```bash
./foghorn-daemon -c example.yaml -v
```

Validate configuration only:
```bash
./foghorn-daemon -c example.yaml --dry-run
```

Verify Docker images are available locally:
```bash
./foghorn-daemon -c example.yaml --verify-image-availability
```

Combine flags for full validation:
```bash
./foghorn-daemon -c example.yaml --dry-run --verify-image-availability
```

Manage secrets:
```bash
export FOGHORN_SECRET_MASTER_KEY="$(openssl rand -base64 32)"
printf '%s' 'smtp-password' | ./foghorn-daemon secret set smtp/password
./foghorn-daemon secret list
./foghorn-daemon secret delete smtp/password
```

The scheduler will load the configuration and execute checks based on their cron schedules.

### TUI Client Options

- `-u, --status-url <url>`: Daemon status API base URL (default: `http://127.0.0.1:7676`)
- `-l, --log-level <level>`: Display label in header

## Docker Check Interface

Foghorn executes Docker containers as health checks and communicates with them through a well-defined interface.

### Environment Variables

Check containers receive the following environment variables from Foghorn:

- `FOGHORN_CHECK_NAME`: Name of the check
- `FOGHORN_CHECK_CONFIG`: JSON string with check-specific configuration (from metadata)
- `FOGHORN_ENDPOINT`: Target endpoint to check (if applicable)
- `FOGHORN_TIMEOUT`: Timeout duration for the check

Secret injection:
- Set config env values to `secret://<key>` (example: `SMTP_PASSWORD: secret://smtp/password`)
- Foghorn resolves secrets only at runtime in memory
- Resolved secrets are written to ephemeral files mounted at `/run/foghorn/secrets`
- For each secret env key `NAME`, Foghorn injects `NAME_FILE=/run/foghorn/secrets/NAME`
- Check containers should read from the `_FILE` path, not the env variable itself
- This approach prevents secrets from appearing in container logs or process listings

Example check container reading a secret:
```bash
# Read secret from the file path (recommended)
PASSWORD=$(cat "$SMTP_PASSWORD_FILE")
# Do NOT read from $SMTP_PASSWORD directly (contains "secret://..." reference)
```

All other non-secret environment variables in config are passed directly to the container.

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

See `containers/disk-check/`, `containers/http-check/`, `containers/openssl-check/`, `containers/mail-send-receive-check/`, and `containers/env-dump-check/` for maintained check container implementations of the Foghorn check interface.
Built-in check containers live under `containers/` and are published to GHCR on change.

## Configuration

Foghorn uses YAML configuration files to define checks. See `example.yaml` for a comprehensive example of all available configuration options.

Configuration includes:
- Global settings (version, max concurrent checks)
- Check definitions with container images
- Schedules (cron expressions or intervals)
- Evaluation rules for results
- Metadata and tags
- Environment variables and timeouts

### Container Image Versions

Check containers must use semantic version tags. Supported selectors are `MAJOR`, `MAJOR.PATCH`, and full `MAJOR.MINOR.PATCH`. Partial selectors resolve to the highest matching local version.

### Global Settings

- `version`: Configuration file version (optional)
- `max_concurrent_checks`: Maximum number of checks that can run simultaneously (optional, defaults to unlimited)
- `state_log_period`: Retention period for state log records (optional, required when state log file is set)
- `state_log_file`: Optional state log file path (CLI `--state-log-file` overrides)
- `secret_store_file`: Optional encrypted secret store file path (CLI `--secret-store-file` overrides)

### Concurrency Control

Foghorn supports limiting concurrent check execution to prevent resource exhaustion:

```yaml
max_concurrent_checks: 5
```

When the concurrency limit is reached:
- Checks are queued until a slot becomes available
- No checks are dropped or lost
- Queue is processed in FIFO order

## Scheduler

The scheduler component manages check execution based on cron expressions:
- Parses standard cron expressions (minute, hour, day, month, day of week)
- Calculates next execution time for each check
- Triggers check execution when scheduled time is reached
- Supports time zones for accurate scheduling
- Only executes enabled checks
