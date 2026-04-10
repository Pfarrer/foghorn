## Context

Check containers live under `containers/<name>/` with a `Dockerfile`, `check.sh`, `VERSION`, and `README.md`. The `container-release.yml` workflow triggers on changes to `containers/**`, detects which folders changed via git diff, and builds/pushes only those containers to GHCR.

## Goals / Non-Goals

**Goals:**
- Document the local container infrastructure and automated release workflow

**Non-Goals:**
- Adding new containers, changing the folder structure, or modifying the workflow

## Decisions

**One subfolder per container**: Each container has its own directory with all necessary files (Dockerfile, check script, VERSION, README).

**Change detection via git diff**: The workflow uses `git diff --name-only` against the base commit to detect changed container folders, then builds only those.

**VERSION file for semver**: Each container has a `VERSION` file containing a `MAJOR.MINOR.PATCH` version. The workflow reads this and tags images accordingly.

## Risks / Trade-offs

- [Migration doc-only] → No functional risk
