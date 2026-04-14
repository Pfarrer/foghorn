## 1. Snapshot Counts

- [ ] 1.1 Add `Error int` field to `SnapshotCounts` struct in `scheduler/snapshot.go`
- [ ] 1.2 Add `case "error"` to the snapshot count switch in `scheduler/snapshot.go` `Snapshot()` method
- [ ] 1.3 Update tests in `scheduler/snapshot_test.go` to cover `error` status counting

## 2. Scheduler Status Handling

- [ ] 2.1 Add `case "error"` to the status count switch in `scheduler/scheduler.go` `GetCounts()` method
- [ ] 2.2 Update tests in `scheduler/scheduler_test.go` to verify `error` status is handled

## 3. One-Shot Mode

- [ ] 3.1 Update the one-shot status switch in `internal/daemon/app.go` to treat `error` as unhealthy (non-zero exit)
- [ ] 3.2 Update tests in `internal/daemon/app_test.go` to cover `error` status in one-shot mode

## 4. TUI Rendering

- [ ] 4.1 Update `statusSymbol()` in `tui/model.go` to render `error` with a distinct symbol and color (e.g., `!` in amber/yellow)
- [ ] 4.2 Add `colorError` lipgloss style in TUI styles (distinct from `colorFail` and `colorWarn`)
- [ ] 4.3 Update `formatHistorySymbols()` in `tui/model.go` to handle `error` status
- [ ] 4.4 Update tests in `tui/model_test.go` to cover `error` rendering

## 5. Verification

- [ ] 5.1 Build: `go build ./...` — no errors or warnings
- [ ] 5.2 Run all tests: `go test ./...` — all pass
