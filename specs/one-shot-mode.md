# One-Shot Mode

## Category
functional

## Description
Add a one-shot run mode that executes each configured check once, evaluates results, reports output, and exits.

## Usage Steps
1. Start Foghorn with one-shot mode enabled.
2. Foghorn loads config and resolves check images as normal.
3. Foghorn runs each configured check one time.
4. Foghorn evaluates and reports results, then exits.

## Implementation Notes
- Add a CLI flag (for example `--one-shot`) to enable one-shot mode.
- In one-shot mode, bypass recurring scheduling and trigger each configured check exactly once.
- Reuse existing check execution, timeout, and evaluation paths.
- Preserve current logging and state updates for each run.
- Exit with code `0` when all checks pass; non-zero when any check fails or cannot run.

## Acceptance Criteria
- [ ] CLI supports enabling one-shot mode.
- [ ] Each configured check runs exactly once in one-shot mode.
- [ ] Process exits after all checks complete.
- [ ] Existing evaluator logic is used for one-shot results.
- [ ] Exit code reflects aggregate run success/failure.
