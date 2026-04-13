## 1. Config Extension

- [ ] 1.1 Add `PersistentMemory bool` field to `CheckConfig` in `config/types.go` with yaml tag `persistent_memory`
- [ ] 1.2 Add unit tests for parsing `persistent_memory` field in `config/config_test.go`

## 2. Executor: Volume Management

- [ ] 2.1 Add volume name sanitization function in `executor/docker.go` (lowercase, non-alphanumeric to `-`, prefix `foghorn-memory-`)
- [ ] 2.2 Add volume creation method on `DockerExecutor` that creates a named volume if it does not exist (using `VolumeCreate` and `VolumeInspect` APIs)
- [ ] 2.3 Add unit tests for volume name sanitization

## 3. Executor: Volume Mounting

- [ ] 2.4 In `Execute`, when `checkConfig.PersistentMemory` is true, create/ensure volume and add bind mount to `HostConfig.Binds` at `/run/foghorn/memory`
- [ ] 2.5 In `buildEnvVars`, inject `FOGHORN_PERSISTENT_DIR=/run/foghorn/memory` when persistent memory is enabled
- [ ] 2.6 Add integration test for persistent memory: verify volume is created and mounted when enabled, and absent when disabled

## 4. Verification

- [ ] 4.1 Run `go build ./...` and confirm no compile errors
- [ ] 4.2 Run `go test ./...` and confirm all tests pass
- [ ] 4.3 Run `go vet ./...` and confirm no warnings
