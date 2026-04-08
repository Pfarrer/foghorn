## ADDED Requirements

### Requirement: Separate binary entrypoints
The repository SHALL have `cmd/foghorn-daemon` and `cmd/foghorn-tui` as separate `main` packages. No runtime source files SHALL exist at the repository root.

#### Scenario: Daemon binary
- **WHEN** `cmd/foghorn-daemon` is built and run
- **THEN** the daemon starts the scheduler, executor, and status API

#### Scenario: TUI binary
- **WHEN** `cmd/foghorn-tui` is built and run
- **THEN** the TUI connects to the daemon status API

#### Scenario: No root-level main.go
- **WHEN** the repository root is inspected
- **THEN** no `main.go` or runtime source files exist at the root

### Requirement: Independent lifecycle
The daemon and TUI SHALL be independently startable and stoppable. Stopping one SHALL NOT affect the other.

#### Scenario: TUI stops, daemon continues
- **WHEN** the TUI exits via Ctrl+C
- **THEN** the daemon continues running checks

#### Scenario: Daemon stops, TUI shows error
- **WHEN** the daemon stops while the TUI is running
- **THEN** the TUI displays a connection error and can reconnect when the daemon restarts

### Requirement: Internal packages
Shared, non-public code SHALL reside in `internal/` packages with clear package boundaries.

#### Scenario: Package boundaries
- **WHEN** the `internal/` directory is inspected
- **THEN** it contains packages for daemon orchestration (`internal/daemon`) and status API (`internal/statusapi`)

### Requirement: Local status API
The daemon SHALL expose a local HTTP status API for clients to query check state. The API SHALL be local-only (loopback TCP).

#### Scenario: API accessible locally
- **WHEN** the daemon is running with `--status-listen 127.0.0.1:7676`
- **THEN** the status API responds to HTTP requests on that address

### Requirement: Multi-binary build and release
The CI/CD pipeline SHALL build and release both daemon and TUI binaries for all target architectures.

#### Scenario: Both binaries released
- **WHEN** the release workflow runs
- **THEN** both `foghorn-daemon` and `foghorn-tui` artifacts are produced for each architecture
