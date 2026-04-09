## ADDED Requirements

### Requirement: Version selector parsing
The system SHALL parse image tags into three selector kinds: `MAJOR` (single number, e.g., `1`), `MAJOR.PATCH` (two numbers, e.g., `1.2`), and `MAJOR.MINOR.PATCH` (three numbers, e.g., `1.2.3`). Invalid formats SHALL produce an error.

#### Scenario: Major selector
- **WHEN** `ParseSelector("1")` is called
- **THEN** it returns a selector with kind `SelectorMajor` and major=1

#### Scenario: Major-patch selector
- **WHEN** `ParseSelector("1.2")` is called
- **THEN** it returns a selector with kind `SelectorMajorPatch`, major=1, patch=2

#### Scenario: Full selector
- **WHEN** `ParseSelector("1.2.3")` is called
- **THEN** it returns a selector with kind `SelectorFull`, major=1, minor=2, patch=3

#### Scenario: Invalid selector
- **WHEN** `ParseSelector("1.2.3.4")` or `ParseSelector("1.x")` is called
- **THEN** an error is returned

### Requirement: Partial selector resolution
For `MAJOR` and `MAJOR.PATCH` selectors, the system SHALL resolve against a list of available versions and return the latest matching version.

#### Scenario: Major selector resolves to latest
- **WHEN** selector is `1` and available versions are `1.0.1`, `1.2.1`, `1.1.2`, `2.0.0`
- **THEN** resolution returns `1.2.1` (highest matching major)

#### Scenario: Major-patch selector resolves to latest
- **WHEN** selector is `1.2` and available versions include `1.1.2`, `1.2.1`
- **THEN** resolution returns `1.1.2` (matching major and patch)

#### Scenario: No matching version
- **WHEN** selector is `3` and no version has major=3
- **THEN** resolution returns false (no match)

### Requirement: Full selector exact match
For `MAJOR.MINOR.PATCH` selectors, the system SHALL use the exact tag without resolution.

#### Scenario: Full selector matches exactly
- **WHEN** selector is `1.2.3`
- **THEN** only `1.2.3` matches; `1.2.4` does not

### Requirement: Image reference parsing
The system SHALL parse image references in the format `repository:tag`. Digest references (`@sha256:...`) and the `latest` tag SHALL be rejected.

#### Scenario: Valid reference
- **WHEN** `ParseReference("foghorn/disk-check:1.0.0")` is called
- **THEN** it returns repository=`foghorn/disk-check`, tag=`1.0.0`

#### Scenario: Digest rejected
- **WHEN** the image string contains `@`
- **THEN** an error is returned

#### Scenario: Latest tag rejected
- **WHEN** the tag is `latest`
- **THEN** an error is returned
