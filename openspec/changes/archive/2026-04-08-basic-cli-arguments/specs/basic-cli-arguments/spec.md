## ADDED Requirements

### Requirement: Help flag
The daemon SHALL accept `-h` and `--help` flags. When set, it SHALL display usage information and exit with code 0.

#### Scenario: Short help flag
- **WHEN** the daemon is invoked with `-h`
- **THEN** usage information is printed to stderr and the process exits with code 0

#### Scenario: Long help flag
- **WHEN** the daemon is invoked with `--help`
- **THEN** usage information is printed to stderr and the process exits with code 0

### Requirement: Config path flag
The daemon SHALL accept `-c` and `--config` flags with a required file path argument. The path is required — omitting it SHALL produce an error and exit with code 1.

#### Scenario: Config path provided
- **WHEN** `-c /path/to/config.yaml` is provided
- **THEN** the daemon loads the configuration from that path

#### Scenario: Config path omitted
- **WHEN** no `-c`/`--config` flag is provided
- **THEN** the daemon prints "configuration file path is required" and exits with code 1

### Requirement: Verbose flag
The daemon SHALL accept `-v` and `--verbose` flags. When set, verbose logging SHALL be enabled.

#### Scenario: Verbose enabled
- **WHEN** `-v` is set
- **THEN** the logger is initialized in verbose mode

#### Scenario: Verbose combined with log level
- **WHEN** both `-v` and `-l debug` are set
- **THEN** verbose mode is enabled alongside the debug log level

### Requirement: Dry-run flag
The daemon SHALL accept `-d` and `--dry-run` flags. When set, the daemon SHALL validate the configuration and exit without starting the scheduler.

#### Scenario: Dry-run with valid config
- **WHEN** `-d` is set and the config is valid
- **THEN** the daemon validates successfully and exits without running checks

#### Scenario: Dry-run with invalid config
- **WHEN** `-d` is set and the config has errors
- **THEN** the daemon reports validation errors and exits with code 1

### Requirement: Log level flag
The daemon SHALL accept `-l` and `--log-level` flags with a string argument (default: `info`). Accepted values SHALL be `debug`, `info`, `warn`, `error`. Invalid values SHALL produce an error.

#### Scenario: Valid log level
- **WHEN** `-l debug` is provided
- **THEN** the logger is initialized at debug level

#### Scenario: Invalid log level
- **WHEN** `-l trace` is provided
- **THEN** the daemon prints an error and exits with code 1

### Requirement: Help text completeness
The help output SHALL list all flags with their short and long forms, argument descriptions, and default values.

#### Scenario: Help includes all flags
- **WHEN** the help output is displayed
- **THEN** it includes `-c/--config`, `-l/--log-level`, `-v/--verbose`, `-d/--dry-run`, `-i/--verify-image-availability`, `--status-listen`, `-s/--state-log-file`, `--secret-store-file`, and `-h/--help`

### Requirement: Flag combinations
All flags SHALL work in combination without conflict.

#### Scenario: Multiple flags together
- **WHEN** `-c config.yaml -v -l debug` is provided
- **THEN** the daemon loads the config, enables verbose mode, and sets debug log level
