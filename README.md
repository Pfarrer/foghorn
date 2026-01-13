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
