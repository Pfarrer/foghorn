## 1. CLI Flag Implementation

- [ ] 1.1 Add `--one-shot` flag in `internal/daemon/app.go` (short flag optional)
- [ ] 1.2 Add flag handling in the flag parsing section

## 2. One-Shot Execution Flow

- [ ] 2.1 Create a `runOneShot(cfg)` function that:
  - Takes the loaded config
  - Iterates over enabled checks
  - Executes each check using the Docker executor directly (bypassing scheduler)
  - Respects `max_concurrent_checks` for parallel execution
- [ ] 2.2 Track aggregate results (pass/warn/fail/unknown counts)
- [ ] 2.3 After all checks complete, log summary and exit

## 3. Integration with Existing Code

- [ ] 3.1 Reuse `DockerExecutor.Execute()` for actual check execution
- [ ] 3.2 Reuse evaluator logic for result determination
- [ ] 3.3 Use existing logging and state recording

## 4. Testing

- [ ] 4.1 Add test: one-shot flag enables mode
- [ ] 4.2 Add test: each enabled check runs exactly once
- [ ] 4.3 Add test: exit code 0 when all pass/warn
- [ ] 4.4 Add test: exit code 1 when any fail/unknown
- [ ] 4.5 Verify code compiles without warnings and all tests pass

## 5. Documentation

- [ ] 5.1 Delete legacy spec file `specs/one-shot-mode.md`
- [ ] 5.2 Update `specs/STATUS.md` to remove the entry

## 6. Archive Change

- [ ] 6.1 Archive the OpenSpec change once implementation is complete
