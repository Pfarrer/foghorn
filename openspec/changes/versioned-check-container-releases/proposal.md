## Why

The versioned check container release system is implemented in `containerimage/selector.go` and `containerimage/reference.go`, with a legacy spec at `specs/versioned-check-container-releases.md`. Migrating to OpenSpec formalizes the version selector and resolution logic.

## What Changes

- Formalize the versioned container release feature as an OpenSpec change
- Port the legacy spec into the OpenSpec spec format
- No functional changes

## Capabilities

### New Capabilities
- `versioned-check-container-releases`: Semantic version selectors for check container images, supporting MAJOR, MAJOR.PATCH (major+patch), and MAJOR.MINOR.PATCH (full) selectors with resolution against available registry tags

### Modified Capabilities
<!-- None -->

## Impact

- `openspec/specs/` — new spec directory
- `specs/versioned-check-container-releases.md` — legacy spec to be removed
- No code changes
