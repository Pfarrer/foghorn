## 1. CLI Flag Implementation

- [x] 1.1 Add `--one-shot` flag in `internal/daemon/app.go` (short flag optional)
- [x] 1.2 Add flag handling in the flag parsing section

## 2. One-Shot Execution Flow

- [x] 2.1 Create a `runOneShot(cfg)` function that:
  - Takes the loaded config
  - Iterates over enabled checks
  - Executes each check using the Docker executor directly (bypassing scheduler)
  - Respects `max_concurrent_checks` for parallel execution
- [x] 2.2 Track aggregate results (pass/warn/fail/unknown counts)
- [x] 2.3 After all checks complete, log summary and exit

## 3. Integration with Existing Code

- [x] 3.1 Reuse `DockerExecutor.Execute()` for actual check execution
- [x] 3.2 Reuse evaluator logic for result determination
- [x] 3.3 Use existing logging and state recording

## 4. Testing

- [x] 4.1 Add test: one-shot flag enables mode
- [x] 4.2 Add test: each enabled check runs exactly once
- [x] 4.3 Add test: exit code 0 when all pass/warn
- [x] 4.4 Add test: exit code 1 when any fail/unknown
- [x] 4.5 Verify code compiles without warnings and all tests pass

## 5. Documentation

- [x] 5.1 Delete legacy spec file `specs/one-shot-mode.md`
- [x] 5.2 Update `specs/STATUS.md` to remove the entry

## 6. Archive Change

- [ ] 6.1 Archive the OpenSpec change once implementation is complete
