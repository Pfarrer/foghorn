## ADDED Requirements

### Requirement: YAML file parsing
The loader SHALL read a YAML file and parse it into a `Config` struct. If the file cannot be opened or contains invalid YAML, an error SHALL be returned.

#### Scenario: Valid YAML file loaded
- **WHEN** a valid YAML config file is provided
- **THEN** a `Config` struct is returned with all parsed fields

#### Scenario: File not found
- **WHEN** the specified file path does not exist
- **THEN** an error is returned indicating the file could not be opened

#### Scenario: Invalid YAML syntax
- **WHEN** the file contains malformed YAML
- **THEN** an error is returned indicating YAML parse failure

### Requirement: Dual format support
The loader SHALL support both document-per-check format (each YAML document with a `name` field is a check) and list format (checks defined under a `checks:` key). Both formats SHALL be usable in the same file.

#### Scenario: Document-per-check format
- **WHEN** the YAML contains multiple `---` separated documents each with `name`, `image`, `schedule`
- **THEN** all checks are loaded into the `Checks` slice

#### Scenario: List format
- **WHEN** the YAML contains a `checks:` key with a list of check objects
- **THEN** all checks in the list are loaded

#### Scenario: Mixed formats in one file
- **WHEN** a file contains both document-per-check and `checks:` list documents
- **THEN** all checks from both formats are combined

### Requirement: Config validation
The loader SHALL validate required fields on each check: `name` (required), `image` (required, must be a parseable image reference), and `schedule` (at least one of `cron` or `interval`, not both). Global settings SHALL also be validated.

#### Scenario: Missing check name
- **WHEN** a check has an empty `name`
- **THEN** an error is returned indicating name is required

#### Scenario: Missing check image
- **WHEN** a check has no `image` field
- **THEN** an error is returned indicating image is required

#### Scenario: Missing schedule
- **WHEN** a check has neither `cron` nor `interval`
- **THEN** an error is returned indicating schedule is required

#### Scenario: Both cron and interval set
- **WHEN** a check has both `cron` and `interval` set
- **THEN** an error is returned indicating only one should be specified

#### Scenario: Invalid image reference
- **WHEN** a check has an invalid image tag format
- **THEN** an error is returned indicating invalid image tag

#### Scenario: Negative max_concurrent_checks
- **WHEN** `max_concurrent_checks` is negative
- **THEN** an error is returned

### Requirement: Config summary
On successful load, a summary SHALL be printed showing version, total checks, enabled/disabled counts, and max concurrent checks setting.

#### Scenario: Summary printed
- **WHEN** config is loaded successfully
- **THEN** version, check count, enabled/disabled counts, and concurrency limit are printed

### Requirement: Single-load model
The configuration file SHALL be read once at startup. No hot-reloading SHALL occur on file changes.

#### Scenario: Config read once
- **WHEN** the daemon starts
- **THEN** the config file is read once and subsequent file changes have no effect
