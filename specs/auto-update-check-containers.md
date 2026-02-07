# Auto-Update Check Containers

## Category
config

## Description
Add a configuration option to periodically fetch the latest check container versions on a schedule.

## Usage Steps
1. Enable auto-update in config and set the update schedule.
2. Start Foghorn.
3. Foghorn periodically fetches the latest container versions.
4. Checks run using the updated containers.

## Implementation Notes
- Add config fields for `auto_update_containers` (bool) and `auto_update_schedule` (cron or interval).
- Reuse the scheduler to trigger update jobs.
- Update should only pull images that are defined in config.
- Log update attempts and outcomes.
- Failures should not stop normal check execution.

## Acceptance Criteria
- [ ] Config supports enabling/disabling auto-update.
- [ ] Config supports scheduling auto-update.
- [ ] Updates pull the latest versions of configured containers.
- [ ] Auto-update failures do not stop check execution.
- [ ] Update attempts are logged.
