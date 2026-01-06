# Agents

Foghorn monitors server services and properties by running Docker containers as checks on a predefined schedule.

Implemented in Go.

## Architecture

- **Scheduler**: Manages check execution schedules
- **Check Runner**: Executes Docker containers and collects results
- **Evaluator**: Parses and evaluates check results

## Check Execution Flow

1. Scheduler triggers a check
2. Check Runner executes the Docker container
3. Container returns JSON/structured output
4. Evaluator processes the result and determines status

## Development

Implement new check containers by:
1. Creating a Docker container that outputs JSON results
2. Adding the check to the schedule configuration
3. Defining evaluation rules for the output
