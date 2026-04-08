## Context

The daemon CLI (`internal/daemon/app.go`) uses Go's standard `flag` package with both short and long flag variants. It also supports a `secret` subcommand for encrypted secret management. Custom `flag.Usage` provides formatted help output.

## Goals / Non-Goals

**Goals:**
- Document existing CLI flags and behavior as formal OpenSpec artifacts
- Preserve all current behavior

**Non-Goals:**
- Modifying flags, adding new ones, or switching flag libraries

## Decisions

**Standard `flag` package**: Uses Go's built-in flag parsing. No external dependencies.

**Custom usage function**: `flag.Usage` is overridden to provide structured help text with grouped flags.

**Verbose overrides log level**: When `-v`/`--verbose` is set, the logger is initialized in verbose mode regardless of the log level setting.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
