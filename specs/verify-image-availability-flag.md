# Verify Image Availability Flag

## Category
functional

## Description
Add a CLI flag that validates all Docker images referenced in the configuration are available in the local Docker system

## Usage Steps
1. User starts Foghorn with the `--verify-image-availability` flag
2. Foghorn loads and validates the configuration as usual
3. After config validation, Foghorn checks each Docker image in the config
4. For each image, Foghorn queries the local Docker daemon to verify it exists
5. If all images are available, Foghorn proceeds normally or exits (if combined with dry-run)
6. If any image is missing, Foghorn displays a clear error message and exits

## Implementation Notes
Add new CLI flag:
- `--verify-image-availability`: Verify all Docker images in config are available locally

Flag behavior:
- Can be used alone or combined with `--dry-run`
- Runs after configuration validation
- For each check in config, validates `image` field
- Uses Docker client to check if image exists locally (avoiding pull attempts)
- Collects all missing images before reporting (rather than failing on first one)

Validation process:
1. Parse configuration normally
2. Iterate through all enabled checks
3. For each check's image, call Docker API to check local availability
4. Store list of missing images with their associated check names
5. If any images missing, display comprehensive error message and exit with status code 1
6. If all images present, display success message and continue (or exit if dry-run)

Error message format:
```
Error: The following Docker images are not available locally:

- foghorn/disk-check:latest (required by: disk-space-check)
- foghorn/http-check:latest (required by: api-health-check, website-check)

Please pull the missing images:
  docker pull foghorn/disk-check:latest
  docker pull foghorn/http-check:latest
```

Success message format:
```
All Docker images validated successfully:
  - foghorn/disk-check:latest ✓
  - foghorn/http-check:latest ✓
  - foghorn/ping-check:latest ✓
```

Docker API usage:
- Use `ImageInspect` or similar API to check if image exists locally
- Do not attempt to pull images (this is validation only)
- Handle Docker daemon connection errors gracefully

## Acceptance Criteria
- [ ] New CLI flag `--verify-image-availability` and `-i` are added
- [ ] Flag can be combined with existing flags (config, log-level, verbose, dry-run)
- [ ] Validation runs after successful configuration validation
- [ ] Checks all images for enabled checks in the config
- [ ] Uses Docker API to verify local image availability (no pull attempts)
- [ ] Reports all missing images in a single error message with check names
- [ ] Displays helpful `docker pull` commands for missing images
- [ ] Exits with status code 1 when images are missing
- [ ] Displays success message when all images are available
- [ ] Proceeds to normal operation when all images present (unless dry-run)
- [ ] Handles Docker daemon connection errors appropriately
- [ ] Works with multiple checks using the same image (groups them in output)
- [ ] Respects verbose flag for more detailed output
- [ ] Examples in documentation

## Example Usage

Verify images before running:
```bash
./foghorn -c config.yaml --verify-image-availability
```

Combine with dry-run for full validation:
```bash
./foghorn -c config.yaml --dry-run --verify-image-availability
```

Verify with verbose output:
```bash
./foghorn -c config.yaml --verify-image-availability --verbose
```

Successful validation output:
```
Configuration loaded successfully
Version: 1.0
Checks: 3
Enabled checks: 3
Max concurrent checks: 5

Validating Docker images...
  - foghorn/ping-check:latest ✓
  - foghorn/http-check:latest ✓
  - foghorn/disk-check:latest ✓

All Docker images validated successfully.

Scheduler started. Press Ctrl+C to stop.
```

Failed validation output:
```
Configuration loaded successfully
Version: 1.0
Checks: 3
Enabled checks: 3

Validating Docker images...
Error: The following Docker images are not available locally:

- foghorn/ping-check:latest (required by: google-ping)
- custom-check:v2.0 (required by: custom-service-check)

Please pull the missing images:
  docker pull foghorn/ping-check:latest
  docker pull custom-check:v2.0
```

## Passes
true
