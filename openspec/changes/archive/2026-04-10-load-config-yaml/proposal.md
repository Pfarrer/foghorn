## Why

The YAML config loader is implemented in `config/loader.go` with a legacy spec at `specs/load-config-yaml.md`. Migrating to OpenSpec formalizes the loading, validation, and merging behavior.

## What Changes

- Formalize the existing YAML config loading feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `load-config-yaml`: Loads, parses, validates, and merges Foghorn check configurations from YAML files

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/load-config-yaml.md` — legacy spec to be removed
- No code changes
