## Context

The codebase is organized with `cmd/foghorn-daemon` (scheduler, executor, config loading, secret store, status API), `cmd/foghorn-tui` (Bubble Tea dashboard), and shared code in `internal/` packages (daemon, statusapi). The two binaries communicate via a local HTTP status API on a configurable TCP port.

## Goals / Non-Goals

**Goals:**
- Document the daemon/client architecture as formal OpenSpec artifacts

**Non-Goals:**
- Adding new binaries, remote access, or web UI

## Decisions

**Separate binaries**: Daemon and TUI are independent processes that can be started/stopped independently.

**HTTP status API**: The daemon exposes a local HTTP API for clients to query check state. Currently loopback TCP only.

**Go-standard layout**: Entrypoints in `cmd/`, shared code in `internal/`, `go.mod` at repo root.

**No legacy in-process mode**: The old `--tui` flag embedding the TUI in the daemon has been removed.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
