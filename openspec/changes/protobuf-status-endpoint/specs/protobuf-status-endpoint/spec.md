## ADDED Requirements

### Requirement: Configurable endpoint
The config SHALL support `status_proto_port` (int) to set the protobuf/gRPC listener port. When set, the endpoint SHALL start alongside the daemon.

#### Scenario: Endpoint enabled
- **WHEN** `status_proto_port` is set to 9090
- **THEN** the gRPC server starts on port 9090

#### Scenario: Endpoint disabled
- **WHEN** `status_proto_port` is not set
- **THEN** no protobuf endpoint is started

### Requirement: Status data coverage
The protobuf endpoint SHALL expose all information available in the existing JSON status API: check list with name, status, schedule, last run time, next run time, and result history.

#### Scenario: Full status data
- **WHEN** a client fetches status from the protobuf endpoint
- **THEN** the response includes all checks with name, status, schedule type, last run, next run, and history

#### Scenario: Real-time updates
- **WHEN** a check completes while a client is connected
- **THEN** the client receives the updated status (via streaming or subsequent fetch)

### Requirement: Versioned protobuf schema
The protobuf schema SHALL be defined in `.proto` files with a versioned Go package path. Field numbers SHALL never be reused to maintain backward compatibility.

#### Scenario: Schema file exists
- **WHEN** the proto directory is inspected
- **THEN** it contains `.proto` files defining status messages and service

#### Scenario: Backward compatibility
- **WHEN** a new optional field is added to the schema
- **THEN** existing clients continue to work without modification

### Requirement: Client connectivity
A client SHALL be able to connect to the protobuf endpoint and fetch status successfully using the generated Go client code.

#### Scenario: Successful fetch
- **WHEN** a client calls the status RPC
- **THEN** it receives a valid protobuf response with all status data
