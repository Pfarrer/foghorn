## 1. Verify Modularization

- [x] 1.1 Verify `cmd/foghorn-daemon` and `cmd/foghorn-tui` are separate main packages
- [x] 1.2 Verify no runtime source at repo root
- [x] 1.3 Verify daemon and TUI have independent lifecycles
- [x] 1.4 Verify `internal/` packages contain shared logic
- [x] 1.5 Verify local status API works on loopback
- [x] 1.6 Verify CI builds both binaries for all architectures

## 2. Documentation

- [x] 2.1 Delete legacy spec file `specs/daemon-client-modularization.md`
- [x] 2.2 Update `specs/STATUS.md` to remove the entry
- [x] 2.3 Verify all tests pass and code compiles without warnings

## 3. Archive Change

- [x] 3.1 Archive the OpenSpec change once verification is complete
