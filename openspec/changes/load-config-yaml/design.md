## Context

The loader (`config/loader.go`) reads YAML using `gopkg.in/yaml.v3` with a streaming decoder. It supports two formats: document-per-check (each YAML document with a `name` field is a check) and `checks:` list under a top-level config document. Multiple document types are merged into a single `Config`. Validation enforces required fields and valid values.

## Goals / Non-Goals

**Goals:**
- Document the existing config loading behavior as formal OpenSpec artifacts

**Non-Goals:**
- Adding hot-reload, new validation rules, or config format changes

## Decisions

**Streaming YAML decoder**: Uses `yaml.NewDecoder` to process multiple YAML documents in a single file separated by `---`.

**Merge strategy**: Global settings from later documents override earlier ones. Check lists are appended.

**Single-load model**: Configuration is read once at startup; no hot-reloading on file changes.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
