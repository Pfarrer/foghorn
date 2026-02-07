# Protobuf Status Endpoint

## Category
integration

## Description
Expose a port that serves protobuf-based status data for building an external UI with all information from the local interface.

## Usage Steps
1. Configure the status endpoint port.
2. Start Foghorn.
3. Connect a client and fetch protobuf status data.
4. Render an external UI with the same information as the local interface.

## Implementation Notes
- Add config fields for `status_port` and `status_proto` definition.
- Expose a protobuf-based endpoint over TCP (or gRPC if preferred).
- Include checks, schedules, last results, and overall status.
- Keep the data model aligned with the local UI needs.
- Provide a versioned schema and backward compatibility policy.

## Acceptance Criteria
- [ ] A configurable port serves protobuf status data.
- [ ] The endpoint includes all information shown in the local UI.
- [ ] The protobuf schema is versioned and documented.
- [ ] A client can connect and fetch status successfully.
