# Local Check Containers

## Category
integration

## Description
Define common check containers in-repo and build/release them automatically when they change.

## Usage Steps
1. Add a new container under the root containers folder.
2. Document it in the container README.
3. Push changes to the repository.
4. The GitHub Action builds and releases the changed container.

## Implementation Notes
- Create a root folder (e.g. `containers/`) with one subfolder per container.
- Each container folder includes a `README.md`.
- Add a GitHub Action that detects container folder changes and builds/releases only those containers.
- Release process should tag/publish images consistently with existing conventions.

## Acceptance Criteria
- [ ] A root containers folder exists with per-container subfolders.
- [ ] Each container subfolder includes a `README.md`.
- [ ] GitHub Action builds and releases a container when its folder changes.
- [ ] Unchanged containers are not rebuilt.
