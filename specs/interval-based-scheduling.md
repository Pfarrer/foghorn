# Interval-Based Scheduling

## Category
functional

## Description
Add support for interval-based scheduling (e.g., "every 5 minutes") as an alternative to cron expressions

## Usage Steps
1. Define checks with interval schedules in the configuration YAML
2. Use human-readable interval syntax (e.g., "5m", "1h", "30s")
3. Start the Foghorn service
4. Scheduler triggers checks at the configured intervals

## Implementation Notes
- Parse interval syntax from configuration
- Support common time units: s (seconds), m (minutes), h (hours), d (days)
- Calculate next execution time based on intervals
- Integrate interval scheduling with existing scheduler
- Support both cron and interval schedules in same configuration

## Acceptance Criteria
- [ ] Scheduler triggers checks based on configured intervals
- [ ] Interval syntax is parsed correctly for all time units
- [ ] Next execution times are calculated accurately
- [ ] Mixed cron and interval schedules work together
- [ ] Invalid interval syntax produces helpful error messages
