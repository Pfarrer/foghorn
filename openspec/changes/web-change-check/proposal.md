## Why

Detecting changes on web pages is a common monitoring need (compliance pages, status dashboards, pricing pages). Current check containers are stateless — they can only verify current state, not detect changes from the last run. A dedicated check container using a real browser can render JavaScript-heavy pages and compare content against previous runs.

## What Changes

- New check container `web-change-check` under `containers/`
- Uses a headless browser (Playwright/Chromium) to render pages fully (including JS)
- Accepts XPath expressions to scope which part of the page to monitor
- Persists a hash/snapshot of the matched content via `FOGHORN_PERSISTENT_DIR`
- On change detection: fails with diff details, saves new state
- On no change: passes
- Requires `persistent_memory: true` in check config

## Capabilities

### New Capabilities
- `web-change-check`: Check container that renders web pages with a headless browser, extracts content via XPath, persists state, and detects changes between runs

### Modified Capabilities

_(none — this is a standalone container using existing interfaces)_

## Impact

- New directory: `containers/web-change-check/` with Dockerfile, check script, README, VERSION
- Docker image: `ghcr.io/pfarrer/foghorn-web-change-check`
- Depends on `check-persistent-memory` feature (Docker volume for state persistence)
- GitHub Actions workflow auto-builds on changes to `containers/web-change-check/`
