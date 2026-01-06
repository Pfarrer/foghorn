# Foghorn

A service monitoring tool that executes arbitrary Docker containers as health checks.

Checks are Docker containers that run on a predefined schedule, perform custom actions, and return results for evaluation. Unlike traditional monitoring that only supports ping or HTTP checks, Foghorn allows any containerized check.

## Features

- Schedule-based execution of Docker containers
- Arbitrary check logic through containerization
- Result evaluation and status reporting
