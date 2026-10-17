# Task Service

A small task-management service where each user owns a private list of tasks.

## Language

### Billing events

**Invoice due event**:
A fact published by the external billing service stating that an invoice belonging to a user must be paid by a due date. It carries no instruction.
_Avoid_: payment command, billing request

**Payment reminder**:
A task created by this service in reaction to an invoice due event, telling its user to pay the invoice by the due date.
_Avoid_: billing task, invoice task

### Pagination

**Task page**:
The subset of a user's tasks returned by one list request.
_Avoid_: chunk, batch, result set

**Page size**:
The maximum number of tasks a task page may contain.
_Avoid_: per_page

**Offset**:
The number of tasks skipped from the start of the ordered list before a task page begins.
_Avoid_: cursor, after, start

**Has more**:
Whether at least one task exists beyond the current task page.
_Avoid_: next cursor, nextToken, total
