# Load Check Config from YAML File

## Category
functional

## Description
Load and parse check configurations from a YAML file, including check definitions, schedules, and evaluation rules

## Usage Steps
1. Copy example YAML configuration file to have a local copy that can be adjusted
2. Add check definitions to local YAML configuration file
3. Start the Foghorn service to load the configuration
4. Foghorn evaluates and validates the configuration
5. If valid, a summary is printed
6. If invalid, Foghorn exits and prints a helpful problem description

## Implementation Notes
- Parse YAML file and validate required fields (check name, container image, schedule)
- Store check configurations in memory for scheduler
- Handle parsing errors and invalid configurations with clear error messages
- No hot-reloading on file changes, the configuration file is only read once on start
- Validate schedule format (cron expressions or intervals)

## Acceptance Criteria
- [ ] YAML file is parsed without errors
- [ ] All required fields are validated
- [ ] Invalid configurations produce helpful error messages
- [ ] Summary of loaded checks is displayed on successful load

## Passes
false
