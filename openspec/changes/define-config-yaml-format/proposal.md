## Why

The YAML config format is defined in `config/types.go` with a legacy spec at `specs/define-config-yaml-format.md`. Migrating to OpenSpec formalizes the config schema.

## What Changes

- Formalize the existing config YAML format as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `define-config-yaml-format`: YAML configuration schema for Foghorn checks, global settings, scheduling, evaluation rules, and runtime options

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/define-config-yaml-format.md` — legacy spec to be removed
- No code changes
