# Task Service

A small task-management service where each user owns a private list of tasks.

## Language

### Event consumption

**Dead-letter queue**:
A parking queue where this service puts an event it could not consume, so nothing vanishes silently and the failure can be inspected.
_Avoid_: discard pile, error queue

**Poison message**:
An event this service can never consume successfully, such as a malformed invoice due event. It goes straight to the dead-letter queue without retrying.
_Avoid_: bad message, invalid event

**Transient failure**:
A failure that may succeed on a later attempt, such as the database being briefly unavailable.
_Avoid_: temporary error, retryable error

**Permanent failure**:
A failure that no retry can fix, so the event is dead-lettered immediately.
_Avoid_: fatal error, hard failure

### Billing events

**Invoice due event**:
A fact published by the external billing service stating that an invoice belonging to a user must be paid by a due date. It carries no instruction.
_Avoid_: payment command, billing request

**Payment reminder**:
A task created by this service in reaction to an invoice due event, telling its user to pay the invoice by the due date.
_Avoid_: billing task, invoice task

### Task events

**Task created event**:
A fact this service publishes stating that a task came into existence for one of its users, carrying the task's identity, description and creation time. It states a fact and names no action.
_Avoid_: task notification, task message

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
