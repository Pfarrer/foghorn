# Versioned Check Container Releases

## Category
functional

## Description
Check container releases must be versioned with semantic versions. Container images and docs must expose a changelog.

## Usage Steps
1. A check is configured with a version selector: `MAJOR`, two-part `MAJOR.PATCH`, or full `MAJOR.MINOR.PATCH`.
2. For partial selectors (`MAJOR` or `MAJOR.PATCH`), the resolver fetches tags and resolves to the latest matching image tag.
3. For full selectors (`MAJOR.MINOR.PATCH`), the resolver uses the exact image tag.
4. The container README includes a changelog for each released version.

## Implementation Notes
- Treat container image tags as semantic versions.
- Accept selectors `1`, `1.2`, and `1.2.3` (two-part is the "major+patch" selector).
- For partial selectors, fetch available tags from the container registry before resolving.
- Resolve partial selectors to the latest matching semantic version.
- Require each container README to include a changelog section or file.

## Acceptance Criteria
- [x] Version selectors support `MAJOR`, two-part `MAJOR.PATCH`, and full `MAJOR.MINOR.PATCH`.
- [x] Partial selector resolution fetches registry tags and selects the latest matching semantic version.
- [x] Full selectors resolve to the exact specified image tag.
- [x] Each check container README includes a changelog for releases.

## Passes
true
