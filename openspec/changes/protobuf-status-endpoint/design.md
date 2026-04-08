## Context

The daemon currently exposes a JSON HTTP status API (`internal/statusapi`) used by the TUI client. This feature adds a protobuf-based endpoint that serves the same data in a binary format for external UIs and programmatic integrations.

## Goals / Non-Goals

**Goals:**
- Define a protobuf schema covering all status data exposed by the JSON API
- Serve the protobuf endpoint on a configurable port
- Maintain data model alignment with the existing JSON API
- Version the protobuf schema

**Non-Goals:**
- Replacing the JSON HTTP API
- Adding authentication or authorization to the endpoint
- Building an external UI client

## Decisions

**gRPC over HTTP/2**: Use gRPC rather than raw protobuf-over-TCP for standardized transport, code generation, and ecosystem tooling.

**Separate port**: The protobuf endpoint runs on its own configurable port, independent of the JSON HTTP status API.

**Schema versioning**: Proto package uses semantic versioning in the Go package path. Fields use protobuf field numbers (never reuse) for backward compatibility.

## Risks / Trade-offs

- [Build complexity] → Protobuf code generation adds a build step (protoc); mitigated by well-documented Go protobuf tooling
- [Port conflicts] → Separate port may conflict; mitigated by configurable port with sensible default
