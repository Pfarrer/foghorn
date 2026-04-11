## Context

The existing state log (`state/log.go`) stores JSON-line records with `check_name`, `status`, `duration_ms`, and `completed_at`. The scheduler keeps an in-memory `CheckHistoryEntry` slice (last 10 entries) per check. The spec requires richer records and query capabilities.

## Goals / Non-Goals

**Goals:**
- Add start time and error details to state log records
- Support querying history by check name and time range
- Enforce configurable retention policy
- Expose history data via the status API

**Non-Goals:**
- Building a separate history database or external storage
- Full-text search on error messages
- Alerting based on history patterns

## Decisions

**Extend existing state log**: Add fields to the `Record` struct rather than creating a new storage system. Old records without new fields default to zero values.

**JSON-lines format**: Keep the existing JSON-lines format, appending new fields as optional JSON properties for backward compatibility.

**Query in memory**: Load records from file and filter in memory. For the expected data volume (thousands of records), this is sufficient.

## Risks / Trade-offs

- [File growth] → More fields per record increase file size; mitigated by retention pruning
- [Backward compat] → New fields in existing format; old readers ignore unknown JSON keys
