## ADDED Requirements

### Requirement: Check configuration fields
Each check SHALL support the following fields: `name` (string, required), `image` (string, required), `schedule` (object with `cron` and/or `interval`), `evaluation` (list of evaluation rules), `description` (string, optional), `tags` (list of strings, optional), `enabled` (bool, default false), `env` (map of string to string, optional), `timeout` (Go duration string, optional), `metadata` (map, optional).

#### Scenario: Minimal check config
- **WHEN** a check has only `name`, `image`, `schedule`, and `enabled: true`
- **THEN** the config is valid and loadable

#### Scenario: Full check config
- **WHEN** a check uses all fields including `env`, `evaluation`, `tags`, `description`, `timeout`
- **THEN** all fields are parsed correctly

### Requirement: Schedule format
The `schedule` field SHALL support `cron` (standard 5-field cron expression) and `interval` (Go duration string like `5m`, `1h30s`). A check SHALL have at least one schedule type defined.

#### Scenario: Cron schedule
- **WHEN** `schedule.cron` is set to `*/5 * * * *`
- **THEN** the check is scheduled with that cron expression

#### Scenario: Interval schedule
- **WHEN** `schedule.interval` is set to `5m`
- **THEN** the check is scheduled with a 5-minute interval

### Requirement: Evaluation rules
Each evaluation rule SHALL support: `type` (string, required), `condition` (string, required), `threshold` (float, optional), `expected` (any, optional), `metadata` (map, optional).

#### Scenario: Evaluation rule with threshold
- **WHEN** an evaluation rule has type, condition, and threshold
- **THEN** the rule is parsed and stored correctly

### Requirement: Global configuration
The top-level config SHALL support: `checks` (list of check configs), `global` (map, optional), `version` (string, optional), `max_concurrent_checks` (int, optional), `state_log_file` (string, optional), `state_log_period` (string, optional), `secret_store_file` (string, optional), `check_container_debug_output` (string, optional), `debug_output_max_chars` (int, optional).

#### Scenario: Global settings parsed
- **WHEN** the YAML includes `max_concurrent_checks`, `state_log_file`, and `state_log_period`
- **THEN** all global settings are populated correctly

### Requirement: Environment variable support
Check configs SHALL support an `env` map. Each key-value pair is injected as an environment variable into the check container. Values MAY reference secrets via Foghorn's secret injection mechanism.

#### Scenario: Env vars injected
- **WHEN** a check has `env: {HOST: "example.com", PORT: "443"}`
- **THEN** those variables are available in the check container
