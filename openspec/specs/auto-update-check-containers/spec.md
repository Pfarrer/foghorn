## ADDED Requirements

### Requirement: Auto-update configuration
The config SHALL support `auto_update_containers` (bool, default: false) to enable/disable the feature and `auto_update_schedule` (schedule object with cron or interval) to control the update frequency.

#### Scenario: Auto-update enabled
- **WHEN** `auto_update_containers` is `true` and `auto_update_schedule` is set to `interval: "6h"`
- **THEN** auto-update runs every 6 hours

#### Scenario: Auto-update disabled
- **WHEN** `auto_update_containers` is `false` or unset
- **THEN** no auto-update runs occur

#### Scenario: Missing schedule
- **WHEN** `auto_update_containers` is `true` but `auto_update_schedule` is not set
- **THEN** config validation fails with an error

### Requirement: Image pull for configured containers
On each auto-update run, the system SHALL iterate all enabled checks and pull the latest image matching each check's image reference (resolving semver selectors via the image resolver).

#### Scenario: All check images pulled
- **WHEN** auto-update runs and 3 checks are enabled
- **THEN** the latest image is pulled for each of the 3 checks

#### Scenario: Semver selector resolved
- **WHEN** a check uses image `foghorn/openssl-check:1`
- **THEN** the resolver fetches registry tags and pulls the latest `1.x.x` image

### Requirement: Non-blocking failure handling
Auto-update failures SHALL be logged as warnings but SHALL NOT stop the scheduler or prevent check execution.

#### Scenario: Pull failure logged
- **WHEN** a Docker pull fails (e.g., network error)
- **THEN** a warning is logged with the check name and error details

#### Scenario: Checks continue after failure
- **WHEN** an auto-update pull fails
- **THEN** previously pulled images continue to be used for check execution

### Requirement: Update logging
Each auto-update attempt and outcome SHALL be logged.

#### Scenario: Successful update logged
- **WHEN** an image pull succeeds
- **THEN** an info message is logged indicating the check name and new image

#### Scenario: Failed update logged
- **WHEN** an image pull fails
- **THEN** a warning is logged with the check name and error

#### Scenario: No-update-needed logged
- **WHEN** the image is already at the latest version
- **THEN** a debug message is logged indicating the image is up to date
