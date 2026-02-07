# Versioned Check Container Releases

## Category
functional

## Description
Check container releases must be versioned with semantic versions. Container images and docs must expose a changelog.

## Usage Steps
1. A check is configured with a version selector: `MAJOR`, two-part `MAJOR.PATCH`, or full `MAJOR.MINOR.PATCH`.
2. The scheduler resolves the selector to an image tag for the check container.
3. The container README includes a changelog for each released version.

## Implementation Notes
- Treat container image tags as semantic versions.
- Accept selectors `1`, `1.2`, and `1.2.3` (two-part is the "major+patch" selector).
- Define deterministic resolution rules for partial selectors.
- Require each container README to include a changelog section or file.

## Acceptance Criteria
- [x] Version selectors support `MAJOR`, two-part `MAJOR.PATCH`, and full `MAJOR.MINOR.PATCH`.
- [x] Image tag resolution for partial selectors is deterministic and documented.
- [x] Each check container README includes a changelog for releases.

## Passes
true
