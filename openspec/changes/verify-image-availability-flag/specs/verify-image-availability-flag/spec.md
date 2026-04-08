## ADDED Requirements

### Requirement: CLI flag registration
The daemon SHALL accept `-i` and `--verify-image-availability` flags. The flag SHALL be combinable with all other flags including `--dry-run`.

#### Scenario: Short flag
- **WHEN** the daemon is invoked with `-i`
- **THEN** image verification is enabled

#### Scenario: Long flag
- **WHEN** the daemon is invoked with `--verify-image-availability`
- **THEN** image verification is enabled

#### Scenario: Combined with dry-run
- **WHEN** both `-i` and `-d` are set
- **THEN** config validation and image verification run, but the scheduler does not start

### Requirement: Image verification execution
When the flag is set, verification SHALL run after successful config validation. For each enabled check, the image reference SHALL be resolved and inspected locally via Docker API. Images SHALL NOT be pulled.

#### Scenario: Verification runs after config load
- **WHEN** the flag is set and config loads successfully
- **THEN** Docker images for all enabled checks are inspected

#### Scenario: Only enabled checks verified
- **WHEN** the config has both enabled and disabled checks
- **THEN** only images for enabled checks are verified

#### Scenario: No pull attempts
- **WHEN** an image is not found locally
- **THEN** a `docker pull` suggestion is printed but no pull is attempted

### Requirement: Missing image reporting
When images are missing, the daemon SHALL report all missing images with their associated check names in a single error message, include `docker pull` commands, and exit with code 1.

#### Scenario: Single missing image
- **WHEN** one image is not available locally
- **THEN** an error message lists the image, the check requiring it, and a `docker pull` command

#### Scenario: Multiple missing images
- **WHEN** multiple images are not available locally
- **THEN** all missing images are listed with their check names and pull commands in one message

#### Scenario: Multiple checks same image
- **WHEN** two checks use the same missing image
- **THEN** the image is listed once with both check names grouped

### Requirement: Successful verification
When all images are available locally, the daemon SHALL print a success message listing each validated image and continue normal operation.

#### Scenario: All images present
- **WHEN** all enabled check images are available locally
- **THEN** a success message lists each image and the scheduler starts normally

### Requirement: Docker daemon errors
If the Docker daemon is unreachable, the daemon SHALL report a connection error and exit with code 1.

#### Scenario: Docker daemon not running
- **WHEN** the Docker daemon cannot be contacted
- **THEN** an error message indicates the daemon connection failed and the process exits with code 1
