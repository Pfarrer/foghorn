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

## Quality Assurance

Test coverage is critical. Before considering any work complete, always verify:
1. Code compiles without errors
2. No compiler warnings are generated
3. All tests pass successfully
4. All documentation and spec files have been updated

## Spec Status Tracking

Spec implementation status is tracked in [specs/STATUS.md](specs/STATUS.md). When implementing a spec:
1. Move the spec from "Ready" to "Done" upon completion
2. Set the `Passes` field to `true` in the spec file
3. Update any relevant documentation

## Documentation Style

All documentation and markdown files must be concise and short. Minimize the amount of text wherever possible. Be direct and avoid unnecessary elaboration.
