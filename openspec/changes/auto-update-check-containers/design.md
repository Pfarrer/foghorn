## Context

Check container images are referenced by semver selectors (e.g., `foghorn/openssl-check:1`). Currently, images must be manually pulled or pre-built. This feature adds a scheduled job that pulls the latest matching image for each configured check container, so checks always run with up-to-date images without manual intervention.

## Goals / Non-Goals

**Goals:**
- Enable/disable auto-update via config
- Schedule auto-update on a cron or interval basis
- Pull latest images for all configured check containers
- Log update attempts and outcomes
- Failures must not stop normal check execution

**Non-Goals:**
- Automatic restart of running checks after update
- Rolling back to previous image versions
- Image verification or signing

## Decisions

**Reuse scheduler**: The auto-update job runs as a scheduled task using the existing scheduler, same as regular checks but with its own handler.

**Docker pull only**: Auto-update runs `docker pull` for each check's image reference. The image resolver (`containerimage` package) resolves semver selectors to the latest matching tag before pulling.

**Non-blocking failures**: Auto-update failures are logged as warnings but do not affect the scheduler or check execution.

## Risks / Trade-offs

- [Pull bandwidth] → Auto-update may consume bandwidth on each pull; mitigated by configurable schedule (not every tick)
- [Image instability] → Latest image may introduce breaking changes; users can pin exact versions to opt out
