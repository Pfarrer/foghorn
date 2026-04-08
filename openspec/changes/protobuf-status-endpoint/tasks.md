## 1. Protobuf Schema

- [ ] 1.1 Create `proto/` directory with `.proto` files defining status messages (CheckStatus, SchedulerState, etc.) and gRPC service
- [ ] 1.2 Add proto field numbers for all status fields matching the JSON API data model
- [ ] 1.3 Set up protobuf code generation (protoc + Go plugins)

## 2. Config Changes

- [ ] 2.1 Add `status_proto_port` (int, omitempty) to `Config` struct in `config/types.go`
- [ ] 2.2 Add validation for valid port range

## 3. Server Implementation

- [ ] 3.1 Create gRPC server package that implements the status service
- [ ] 3.2 Bridge the gRPC service to the scheduler's `GetAllChecks()` and `GetCheckStatus()` methods
- [ ] 3.3 Support both unary RPC (single fetch) and server streaming (real-time updates)

## 4. Daemon Integration

- [ ] 4.1 Start the gRPC server in `internal/daemon/app.go` when `status_proto_port` is configured
- [ ] 4.2 Gracefully shut down the gRPC server on daemon stop

## 5. Testing

- [ ] 5.1 Add unit tests for the gRPC service handler
- [ ] 5.2 Add integration test: start server, connect client, fetch status, verify response
- [ ] 5.3 Verify code compiles without warnings and all tests pass

## 6. Documentation

- [ ] 6.1 Delete legacy spec file `specs/protobuf-status-endpoint.md`
- [ ] 6.2 Update `specs/STATUS.md` to remove the entry

## 7. Archive Change

- [ ] 7.1 Archive the OpenSpec change once implementation is complete
