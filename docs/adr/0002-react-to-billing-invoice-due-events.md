# React to billing invoice-due events with payment reminder tasks

Status: accepted

The external billing service publishes `InvoiceDue` facts on the `billing` exchange (routing key `invoice-due`, queue `todo-service.events.billing.invoice-due`). This service subscribes and reacts by creating a payment reminder task for the user the invoice belongs to. The event states a fact and never names the action; deciding whether and what task to create is ours.

## Context

An earlier framing — "create a reminder task for user X" — put the user id and the action in the payload, making the message a command in event clothing: the publisher would effectively decide our behavior. Candidate events whose producer and consumer were the same service (user registered/deleted) were also rejected: a service reacting to its own announcements is just talking to itself.

## Considered options

- **Imperative payload (userId + prescribed action).** Rejected: turns the event into a remote command.
- **UserRegistered / UserDeleted.** Rejected: self-reaction; this service owns users.
- **CalendarEventCreated / EmailReceived.** Viable, deferred.
- **InvoiceDue.** Chosen: producer is genuinely external, and the reaction reuses the existing task-creation path with no new domain concepts.

## Consequences

- The billing service re-publishes due reminders until paid, so replayed events create duplicate reminder tasks. Accepted for v1; no dedup key, no task status.
- Events referencing a user this service does not know are skipped with a clear log (the `task.user_id` foreign key would otherwise fail with a cryptic error).
- The consumer is auto-ack: a handler failure loses the message. Accepted for reminder semantics.
- The payload (`invoiceId`, `userId`, `dueDate`, optional `issuer`/`amount`/`currency`) is a cross-service contract owned jointly with the billing service.

Update: the unknown-user and auto-ack consequences were superseded by [ADR-0003](0003-dead-letter-failed-event-consumption.md) — consumers now run manual ack, and unconsumable events (including events for unknown users) are dead-lettered rather than skipped or lost. The payload contract is unchanged.
