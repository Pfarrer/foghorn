## ADDED Requirements

### Requirement: Log levels
The logger SHALL support four levels: `debug`, `info`, `warn`, `error`. Messages below the configured level SHALL be suppressed.

#### Scenario: Debug level shows all messages
- **WHEN** the log level is set to `debug`
- **THEN** debug, info, warn, and error messages are all output

#### Scenario: Info level suppresses debug
- **WHEN** the log level is set to `info`
- **THEN** info, warn, and error messages are output; debug messages are suppressed

#### Scenario: Error level only
- **WHEN** the log level is set to `error`
- **THEN** only error messages are output

### Requirement: Level parsing
`ParseLevel` SHALL accept string values `debug`, `info`, `warn`, `error`. Invalid values SHALL return an error.

#### Scenario: Valid level strings
- **WHEN** `ParseLevel("warn")` is called
- **THEN** it returns `LevelWarn`

#### Scenario: Invalid level string
- **WHEN** `ParseLevel("trace")` is called
- **THEN** an error is returned

### Requirement: Output format
Log output SHALL use the format `[LEVEL] message` written to stdout. Each level SHALL be displayed as its uppercase string.

#### Scenario: Format
- **WHEN** `logger.Info("check started")` is called
- **THEN** output is `[INFO] check started`

### Requirement: Verbose timestamps
When verbose mode is enabled, each log line SHALL be prefixed with a UTC ISO 8601 timestamp.

#### Scenario: Verbose output
- **WHEN** verbose mode is on and `logger.Info("test")` is called
- **THEN** output starts with a timestamp like `2026-01-18T12:00:00Z [INFO] test`

#### Scenario: Non-verbose output
- **WHEN** verbose mode is off and `logger.Info("test")` is called
- **THEN** output is `[INFO] test` without a timestamp

### Requirement: Global logger
Package-level functions `Debug`, `Info`, `Warn`, `Error` SHALL use a global logger instance. `SetGlobal` SHALL set the global instance; `GetGlobal` SHALL return it (creating a default if nil).

#### Scenario: Global functions work
- **WHEN** `SetGlobal(logger.New(LevelInfo, false))` is called
- **THEN** `logger.Info(...)` writes to stdout at info level

#### Scenario: Default global logger
- **WHEN** `logger.GetGlobal()` is called without `SetGlobal`
- **THEN** a default logger at info level is returned
