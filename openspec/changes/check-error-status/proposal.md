## Why

The check interface currently has no way to distinguish a technical error in the check container itself (network failure, browser crash, timeout) from an actual problem with the monitored property. Containers must use `fail` for both, making it impossible for Foghorn core and consumers to tell whether the check result is trustworthy or the check simply couldn't complete.

## What Changes

- **BREAKING**: Add `error` as a new status value to the check output contract, alongside `pass`, `fail`, `warn`, and `unknown`
- Update the executor to recognize and record `error` status
- Update the TUI dashboard to display `error` status distinctly from `fail`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `docker-check-interface`: Add `error` status value to the JSON output format and define its semantics (check could not complete due to a technical error in the container, not a problem with the checked property)

## Impact

- `docker-check-interface` spec: status enum expanded
- Executor (`executor/`): must parse and store `error` status
- TUI (`cmd/foghorn-tui/`): must render `error` status distinctly
- State log (`state/`): must persist `error` results
- Existing check containers: unaffected (they can continue using only `pass`/`fail`/`warn`/`unknown`)
