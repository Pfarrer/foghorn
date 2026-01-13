# Example Docker Check

This is an example Docker container that implements the Foghorn check interface.

## Features

- Performs HTTP health checks on specified endpoints
- Outputs results in JSON format to stdout
- Returns appropriate exit codes
- Supports environment variables for configuration

## Environment Variables

The check container receives the following environment variables from Foghorn:

- `FOGHORN_CHECK_NAME`: Name of the check
- `FOGHORN_ENDPOINT`: Target endpoint to check (if applicable)
- `FOGHORN_TIMEOUT`: Timeout duration for the check
- `FOGHORN_SECRETS`: JSON string with secrets (API keys, tokens)
- `FOGHORN_CHECK_CONFIG`: JSON string with check-specific configuration

Additional custom environment variables from the check configuration are also available.

## Output Format

The check outputs JSON to stdout with the following structure:

```json
{
  "status": "pass|fail|warn|unknown",
  "message": "Human-readable result description",
  "data": {
    "http_code": 200,
    "endpoint": "https://example.com/health"
  },
  "timestamp": "2025-01-13T12:00:00Z",
  "duration_ms": 150
}
```

## Exit Codes

- `0`: Check passed successfully
- `non-zero`: Check failed or encountered an error

## Building the Example

```bash
cd examples/docker-check
docker build -t foghorn/http-check-example .
```

## Using in Foghorn

Add to your Foghorn configuration:

```yaml
name: "http-health-check"
description: "HTTP health check example"
tags:
  - "network"
  - "health"
enabled: true
image: "foghorn/http-check-example:latest"
schedule:
  cron: "*/5 * * * *"
env:
  ENDPOINT: "https://api.example.com/health"
timeout: "30s"
```

## Customizing

To create your own check container:

1. Create a Dockerfile based on your preferred base image
2. Add your check logic (script, compiled binary, etc.)
3. Ensure your check outputs JSON in the expected format to stdout
4. Return appropriate exit codes
5. Build and push your Docker image
6. Reference it in your Foghorn configuration
