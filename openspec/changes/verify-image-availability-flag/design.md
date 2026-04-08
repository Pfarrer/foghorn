## Context

The `-i`/`--verify-image-availability` flag is defined in `internal/daemon/app.go`. When set, after config validation it iterates enabled checks, resolves image references via `imageresolver`, inspects local Docker images via `ImageInspectWithRaw`, and reports missing images with `docker pull` suggestions.

## Goals / Non-Goals

**Goals:**
- Document the existing image verification feature as formal OpenSpec artifacts

**Non-Goals:**
- Adding image pulling, remote registry verification, or new flags

## Decisions

**Docker API inspect-only**: Uses `ImageInspectWithRaw` to check local availability without pulling. `IsErrNotFound` distinguishes missing images from other errors.

**Image resolution before check**: Uses `imageresolver.Resolve` to resolve semver tags to concrete digests before inspecting, so checks using the same semver tag are grouped.

**Collect-all-errors pattern**: All missing images are collected before reporting, rather than failing on the first missing image.

**Combinable with dry-run**: Works standalone or combined with `--dry-run` for full pre-flight validation.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
