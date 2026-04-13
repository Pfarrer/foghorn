## 1. Config Changes

- [x] 1.1 Add `auto_update_containers` (bool, omitempty) and `auto_update_schedule` (Schedule, omitempty) to `Config` struct in `config/types.go`
- [x] 1.2 Add validation: when `auto_update_containers` is true, `auto_update_schedule` must be set
- [x] 1.3 Update config loader to merge and validate new fields

## 2. Auto-Update Implementation

- [x] 2.1 Create auto-update handler function that iterates enabled checks and pulls latest images
- [x] 2.2 Integrate image resolver to resolve semver selectors before pulling
- [x] 2.3 Log success, failure, and no-update-needed outcomes
- [x] 2.4 Ensure pull failures do not affect scheduler or check execution

## 3. Daemon Integration

- [x] 3.1 Register auto-update as a scheduled job in the daemon when config enables it
- [x] 3.2 Pass the Docker client and image resolver to the auto-update handler

## 4. Testing

- [x] 4.1 Add unit tests for auto-update config validation
- [x] 4.2 Add unit tests for auto-update handler (success, failure, already up-to-date)
- [x] 4.3 Verify code compiles without warnings and all tests pass

## 5. Documentation

- [x] 5.1 Delete legacy spec file `specs/auto-update-check-containers.md`
- [x] 5.2 Update `specs/STATUS.md` to remove the entry

## 6. Archive Change

- [x] 6.1 Archive the OpenSpec change once implementation is complete
