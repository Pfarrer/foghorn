## Context

The TUI client (`cmd/foghorn-tui/main.go`) connects to the daemon's status API and renders a dashboard using Bubble Tea (`github.com/charmbracelet/bubbletea`). The model (`tui/model.go`) handles the display layout, and `tui/remote_status.go` fetches data from the API. Styles are defined in `tui/styles.go`.

## Goals / Non-Goals

**Goals:**
- Document the existing TUI dashboard behavior as formal OpenSpec artifacts

**Non-Goals:**
- Adding interactivity, new sections, or changing the display format

## Decisions

**Bubble Tea**: Uses the Charm Bubble Tea TUI framework for terminal rendering.

**Remote model**: Fetches status from the daemon's HTTP API rather than sharing in-memory state, allowing the TUI to run as a separate process.

**Read-only display**: No user interaction beyond Ctrl+C to exit.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
