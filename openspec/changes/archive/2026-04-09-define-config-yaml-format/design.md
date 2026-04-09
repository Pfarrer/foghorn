## Context

The config schema is defined in `config/types.go` using Go structs with YAML tags. The loader (`config/loader.go`) supports document-per-check and `checks:` list formats. Top-level `Config` struct holds global settings and the check list.

## Goals / Non-Goals

**Goals:**
- Document the existing YAML config format as formal OpenSpec artifacts

**Non-Goals:**
- Adding new config fields or changing the schema

## Decisions

**YAML struct tags**: Config types use `yaml` struct tags directly for marshaling/unmarshaling via `gopkg.in/yaml.v3`.

**Dual format support**: Checks can be defined as YAML documents separated by `---` or as a list under `checks:`.

**Optional fields**: Most config fields use `omitempty` — unset fields take zero values or defaults at runtime.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
