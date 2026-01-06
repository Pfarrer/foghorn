# Basic Scheduler with Cron Scheduling

## Category
functional

## Description
Implement core scheduler component that triggers Docker container checks based on cron expressions

## Usage Steps
1. Define checks with cron schedules in the configuration YAML
2. Start the Foghorn service
3. Scheduler automatically triggers checks based on their cron schedules
4. Monitor check execution and results

## Implementation Notes
- Implement scheduler component to manage check execution schedules
- Support standard cron expressions (minute, hour, day, month, day of week)
- Parse cron expressions from configuration
- Calculate next execution time for each check
- Trigger check execution when scheduled time is reached
- Handle time zones in cron expressions

## Acceptance Criteria
- [ ] Scheduler triggers checks based on configured cron schedules
- [ ] Cron expressions are parsed correctly
- [ ] Next execution times are calculated accurately
- [ ] Checks execute at the correct scheduled times
- [ ] Time zone support is working

## Passes
false
