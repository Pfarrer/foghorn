## Why

The protobuf status endpoint is spec'd but not yet implemented. Currently the daemon exposes a JSON-based HTTP status API used by the TUI client. A protobuf endpoint would enable efficient binary status streaming for external UIs and integrations. Creating this as an OpenSpec change defines the requirements and implementation plan.

## What Changes

- Define the protobuf status endpoint with full OpenSpec artifacts
- Add protobuf schema definitions for check status, scheduler state, and results
- Expose a gRPC or protobuf-over-TCP endpoint alongside the existing JSON HTTP API
- Provide a versioned schema with backward compatibility

## Capabilities

### New Capabilities
- `protobuf-status-endpoint`: Binary protobuf-based status API exposing check states, scheduler info, and results for external client consumption

### Modified Capabilities
<!-- None -->

## Impact

- `config/types.go` — new config fields for protobuf endpoint
- New `proto/` directory with `.proto` schema files
- New package for gRPC/protobuf server implementation
- `internal/daemon/app.go` — start protobuf endpoint alongside JSON API
- `specs/protobuf-status-endpoint.md` — legacy spec to be removed
