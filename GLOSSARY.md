# Task Service

A small task-management service where each user owns a private list of tasks.

## Language

### Pagination

**Task page**:
The subset of a user's tasks returned by one list request.
_Avoid_: chunk, batch, result set

**Page size**:
The maximum number of tasks a task page may contain.
_Avoid_: limit, per_page, pageSize

**Offset**:
The number of tasks skipped from the start of the ordered list before a task page begins.
_Avoid_: cursor, after, start

**Has more**:
Whether at least one task exists beyond the current task page.
_Avoid_: next cursor, nextToken, total
