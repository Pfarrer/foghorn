## Context

The version selector system (`containerimage/selector.go`) parses image tags into semver selectors with three kinds: `SelectorMajor` (e.g., `1`), `SelectorMajorPatch` (e.g., `1.2`), and `SelectorFull` (e.g., `1.2.3`). Partial selectors are resolved against available registry tags to find the latest matching version. Image references (`containerimage/reference.go`) parse `repository:tag` strings.

## Goals / Non-Goals

**Goals:**
- Document the version selector and resolution behavior

**Non-Goals:**
- Adding pre-release/build metadata support or changing selector semantics

## Decisions

**Three selector kinds**: `MAJOR` matches any version with that major number. `MAJOR.PATCH` matches any version with that major and patch (variable minor). `MAJOR.MINOR.PATCH` matches exactly one version.

**Latest-match resolution**: `ResolveSelector` picks the highest matching version from a list using semver comparison.

**No `latest` tag allowed**: The reference parser rejects `latest` to force explicit versioning.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
