## Context

The logger (`logger/logger.go`) is a custom lightweight implementation with four log levels (debug, info, warn, error). It uses a global singleton pattern. Output goes to stdout. Verbose mode adds UTC timestamps.

## Goals / Non-Goals

**Goals:**
- Document the existing logging behavior as formal OpenSpec artifacts

**Non-Goals:**
- Adding structured logging fields, log file output, or external logging libraries

## Decisions

**Custom logger over stdlib**: A simple `fmt.Fprintf`-based logger rather than `log/slog` or a third-party library. Provides `[LEVEL] message` format.

**Global singleton**: `SetGlobal`/`GetGlobal` provides package-level `Debug`, `Info`, `Warn`, `Error` functions.

**Verbose = timestamps**: When verbose mode is enabled, each log line is prefixed with a UTC ISO 8601 timestamp.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
