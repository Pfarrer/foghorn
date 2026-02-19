# Daemon and Client Modularization

## Category
functional

## Description
Split Foghorn into a daemon core and separately runnable client modules. Move runnable source code out of the repository root into Go-standard subfolders.

## Usage Steps
1. Start the daemon process.
2. Connect with the TUI client process.
3. View daemon status and check results from the TUI.
4. Stop and restart either process independently.

## Implementation Notes
- Move entrypoints to `cmd/`:
  - `cmd/foghorn-daemon` for scheduler/executor/evaluator runtime.
  - `cmd/foghorn-tui` for the terminal UI client.
- Move shared, non-public app code to `internal/` packages and keep `go.mod` at repo root.
- Define a local daemon API boundary so clients can query state and stream updates.
- Keep transport local-first for now (loopback TCP or Unix socket); no remote auth scope in this spec.
- Reuse and extend existing status endpoint work so the TUI reads daemon state over IPC instead of in-process calls.
- Treat a future web UI as out of scope for this spec, but keep the daemon API client-friendly.
- Update Docker, CI, and release workflows for multiple binaries.
- Update docs and examples to show separate daemon and TUI startup commands.

## Acceptance Criteria
- [x] Repository root no longer contains tool runtime source files like `main.go`.
- [x] Daemon starts from `cmd/foghorn-daemon` and runs checks on schedule.
- [x] TUI starts from `cmd/foghorn-tui` and reads live state from daemon API.
- [x] Daemon and TUI can be started/stopped independently.
- [x] Shared logic is in `internal/` packages with clear package boundaries.
- [x] Existing behavior for scheduling and check execution remains correct after refactor.
- [x] Build/test/release automation supports multiple binaries.
- [x] Daemon exposes a local status API boundary for clients (Unix socket or loopback TCP).
- [x] Standalone TUI client connects to daemon API and no longer depends on in-process scheduler access.
- [x] Legacy in-process `--tui` mode is removed after standalone client is stable.

Passes: true
