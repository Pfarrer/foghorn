# Define Config YAML Format and Create Example

## Category
config

## Description
Define the YAML configuration format for Foghorn check definitions and create an example configuration file

## Usage Steps
1. Reference the example configuration file for YAML structure
2. Copy and modify the example to define custom checks

## Implementation Notes
- Define YAML schema for check configurations
- Include fields: check name, container image, schedule (cron or interval), evaluation rules
- Support metadata fields (description, tags, enabled status)
- Define evaluation rule format (thresholds, conditions, expected values)
- Create example.yaml with comprehensive examples of all fields
- Document YAML structure in comments within the example file
- Support environment variable substitution for sensitive values

## Acceptance Criteria
- [ ] YAML schema document exists with all fields defined
- [ ] example.yaml file exists in the repository
- [ ] Example includes at least 3 different check types
- [ ] All possible fields are demonstrated in the example
- [ ] YAML is valid and can be parsed without errors
- [ ] Example has clear inline comments explaining each field
- [ ] Schedule formats (cron and interval) are demonstrated

## Passes
true
