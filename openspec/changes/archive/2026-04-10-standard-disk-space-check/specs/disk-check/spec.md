## ADDED Requirements

### Requirement: Mount point validation
The container SHALL validate that `MOUNT_POINT` is provided and exists on the filesystem. Missing or non-existent mount points SHALL produce an `unknown` status.

#### Scenario: Missing MOUNT_POINT
- **WHEN** the `MOUNT_POINT` environment variable is not set
- **THEN** the container outputs JSON with status `unknown` and message "MOUNT_POINT is required"

#### Scenario: Non-existent mount point
- **WHEN** `MOUNT_POINT` is set to a path that does not exist
- **THEN** the container outputs JSON with status `unknown` and message indicating the mount point does not exist

### Requirement: Disk space statistics
The container SHALL read disk space statistics using `df` and report `total_bytes`, `used_bytes`, `free_bytes`, and `usage_percent` in the output data.

#### Scenario: Successful disk stats read
- **WHEN** `MOUNT_POINT` is set to `/host` (a valid mount)
- **THEN** the output data includes `total_bytes`, `used_bytes`, `free_bytes`, and `usage_percent` as numeric values

#### Scenario: Failed disk stats read
- **WHEN** `df` fails to produce output for the given mount point
- **THEN** the container outputs JSON with status `unknown` and message indicating failure to read disk usage

### Requirement: Percentage-based thresholds
The container SHALL support `WARNING_THRESHOLD_PERCENT` (default: 80) and `CRITICAL_THRESHOLD_PERCENT` (default: 90). Usage at or above critical → `fail`; at or above warning (but below critical) → `warn`; below warning → `pass`. Critical MUST be >= warning or the check SHALL produce `unknown`.

#### Scenario: Usage below warning threshold
- **WHEN** usage is 75% and warning threshold is 80%
- **THEN** status is `pass`

#### Scenario: Usage at or above warning threshold
- **WHEN** usage is 85% and warning threshold is 80%, critical threshold is 90%
- **THEN** status is `warn`

#### Scenario: Usage at or above critical threshold
- **WHEN** usage is 92% and critical threshold is 90%
- **THEN** status is `fail`

#### Scenario: Invalid threshold relationship
- **WHEN** `CRITICAL_THRESHOLD_PERCENT` is less than `WARNING_THRESHOLD_PERCENT`
- **THEN** the container outputs status `unknown` with message indicating invalid thresholds

### Requirement: Byte-based thresholds
The container SHALL support optional `WARNING_THRESHOLD_BYTES` and `CRITICAL_THRESHOLD_BYTES`. When set, used bytes at or above critical → `fail`; at or above warning (but below critical) → `warn`. Byte thresholds are evaluated alongside percentage thresholds; the most severe status wins.

#### Scenario: Byte threshold triggers warn
- **WHEN** percentage usage is below warning but used bytes exceeds `WARNING_THRESHOLD_BYTES`
- **THEN** status is `warn`

#### Scenario: Byte threshold triggers fail
- **WHEN** used bytes exceeds `CRITICAL_THRESHOLD_BYTES`
- **THEN** status is `fail`

#### Scenario: Byte thresholds not set
- **WHEN** `WARNING_THRESHOLD_BYTES` and `CRITICAL_THRESHOLD_BYTES` are not set
- **THEN** only percentage thresholds are evaluated

### Requirement: Inode usage checking
The container SHALL support optional inode checking via `CHECK_INODES` (default: `true`). When enabled, inode usage SHALL be read via `df -Pi` and compared against `INODE_WARNING_PERCENT` (default: 85) and `INODE_CRITICAL_PERCENT` (default: 95).

#### Scenario: Inode usage below thresholds
- **WHEN** `CHECK_INODES` is `true` and inode usage is below inode warning threshold
- **THEN** inode check passes and inode stats are included in data

#### Scenario: Inode usage at or above warning
- **WHEN** `CHECK_INODES` is `true` and inode usage is at or above `INODE_WARNING_PERCENT` but below `INODE_CRITICAL_PERCENT`
- **THEN** status is `warn` (unless already `fail` from disk thresholds)

#### Scenario: Inode usage at or above critical
- **WHEN** `CHECK_INODES` is `true` and inode usage is at or above `INODE_CRITICAL_PERCENT`
- **THEN** status is `fail`

#### Scenario: Inode checking disabled
- **WHEN** `CHECK_INODES` is `false`
- **THEN** inode stats are reported as zero and message notes "inode check disabled"

#### Scenario: Inode stats unavailable
- **WHEN** `CHECK_INODES` is `true` but the filesystem reports `-` for inode values
- **THEN** inode checking is gracefully disabled and message notes "inode stats unavailable"

### Requirement: Timeout handling
The container SHALL respect `TIMEOUT_SECONDS` (default: 10) and `FOGHORN_TIMEOUT` for operation timeouts. If `FOGHORN_TIMEOUT` is set and is shorter than `TIMEOUT_SECONDS`, it SHALL take precedence.

#### Scenario: FOGHORN_TIMEOUT overrides
- **WHEN** `FOGHORN_TIMEOUT` is set to a value shorter than `TIMEOUT_SECONDS`
- **THEN** the effective timeout is `FOGHORN_TIMEOUT`

#### Scenario: Invalid timeout value
- **WHEN** `TIMEOUT_SECONDS` is set to a non-parseable value
- **THEN** the container outputs status `unknown` with message indicating invalid timeout

### Requirement: Standard JSON output format
The container SHALL output JSON to stdout conforming to the Foghorn check contract: `{status, message, data, timestamp, duration_ms}`. The `data` object SHALL include `mount_point`, `total_bytes`, `used_bytes`, `free_bytes`, `usage_percent`, `total_inodes`, `used_inodes`, `free_inodes`, and `inode_usage_percent`.

#### Scenario: Valid JSON output on success
- **WHEN** the check completes successfully
- **THEN** stdout is valid JSON with all required fields populated

#### Scenario: Valid JSON output on failure
- **WHEN** the check encounters an error or threshold breach
- **THEN** stdout is valid JSON with appropriate status and descriptive message

### Requirement: Exit code semantics
The container SHALL exit with code 0 for `pass` status and non-zero for `warn`, `fail`, or `unknown` status.

#### Scenario: Pass exit code
- **WHEN** status is `pass`
- **THEN** container exits with code 0

#### Scenario: Non-pass exit code
- **WHEN** status is `warn`, `fail`, or `unknown`
- **THEN** container exits with code 1
