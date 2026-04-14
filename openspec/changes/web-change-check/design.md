## Context

Foghorn check containers are stateless Docker workloads outputting JSON. The new `check-persistent-memory` feature (pending) adds Docker volume persistence. This change builds on that by creating a new check container that leverages persistent memory to detect web page changes across runs.

Existing containers use shell scripts (bash/curl). This container needs a headless browser for JavaScript rendering, which requires a different approach.

## Goals / Non-Goals

**Goals:**
- Render web pages with a real browser (JS execution, dynamic content)
- Extract page content using XPath expressions
- Persist content hash between runs via `FOGHORN_PERSISTENT_DIR`
- Detect changes: fail on change (with diff), pass when unchanged
- Follow existing container conventions (Dockerfile, VERSION, README, JSON output)

**Non-Goals:**
- Visual regression / screenshot comparison
- Authentication flows (login forms, cookies)
- Monitoring response time or status codes (http-check already does this)
- Alerting or notification (handled by Foghorn core)

## Decisions

### 1. Use Python + Playwright

**Decision**: Implement in Python using Playwright with Chromium.

**Rationale**: Playwright provides a clean API for headless browsing with full XPath support. Python has strong HTML/text processing libraries. Playwright's Docker images are well-maintained and reasonably sized.

**Alternatives**:
- **Puppeteer (Node.js)**: Equally capable but Playwright has better multi-browser support and a more stable API.
- **Selenium**: Heavier, more configuration, less reliable for CI containers.
- **Go + chromedp**: Keeps the Go stack, but chromedp API is lower-level and XPath extraction is more verbose.

### 2. Content comparison via SHA-256 hash

**Decision**: Hash the extracted text content with SHA-256 and store the hash in a file at `FOGHORN_PERSISTENT_DIR/state.json`. On change, store both old and new hash plus a snippet of the diff.

**Rationale**: Deterministic, fast comparison. Small storage footprint. The hash file doubles as a simple state marker.

**Alternatives**:
- **Full content snapshot**: More disk usage, harder to compare.
- **Structural comparison (DOM diff)**: Overkill for change detection, complex to implement.

### 3. XPath as the extraction mechanism

**Decision**: Accept `XPATH` env var. When set, extract only the text content of matching nodes. When unset, extract the full page body text.

**Rationale**: XPath is a standard, expressive selection language. Allows scoping to specific sections (e.g., `//div[@id='content']`). Defaulting to full body text covers simple use cases.

### 4. Change detection output format

**Decision**: On change, output `status: "fail"` with a message describing the change. Include `data.previous_hash`, `data.current_hash`, and `data.changed: true` in the output. On no change, output `status: "pass"` with `data.changed: false`.

**Rationale**: Maps naturally to Foghorn's pass/fail semantics. The `data` field provides machine-readable change details.

### 6. Technical errors produce `unknown` status

**Decision**: When the check cannot be performed due to technical errors (network failures, DNS errors, timeouts, browser crashes, SSL failures), output `status: "unknown"` with `data.error_type` categorizing the error.

**Rationale**: Distinguishes "the check could not run at all" from "a change was detected". Without this distinction, a network outage would be indistinguishable from a page change. The `unknown` status is part of the Foghorn check contract (`pass`, `fail`, `warn`, `unknown`) and signals that no meaningful comparison could be made. The `error_type` field (`network`, `timeout`, `ssl`, `browser`) allows downstream consumers to react differently to different error categories.

### 5. Wait strategy for dynamic content

**Decision**: Support `WAIT_SECONDS` env var (default 0) to wait after page load before extracting content. Also support `WAIT_FOR_SELECTOR` to wait until a specific CSS selector appears.

**Rationale**: SPAs and dynamic pages need time to render. A simple time-based wait covers most cases; selector-based wait handles known dynamic elements.

## Risks / Trade-offs

- **Large Docker image** (Playwright + Chromium ~500MB) → Acceptable for monitoring; pulled once, reused. Build with `mcr.microsoft.com/playwright/python` base image.
- **Page rendering flakiness** (dynamic ads, A/B tests) → Mitigation: XPath scoping narrows the watched content. Users should target stable elements.
- **Depends on `check-persistent-memory` feature** → This container requires `persistent_memory: true` in config. Document clearly in README.
- **No authentication support** → Initial version assumes public pages. Can be added later via cookie/header injection.
