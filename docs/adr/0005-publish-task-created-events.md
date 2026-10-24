# Publish task-created events

Status: accepted

The task service publishes a `TaskCreated` fact for every task that comes into existence — created through the HTTP API or as a payment reminder reacting to an invoice-due event alike — on a direct exchange named after the domain (`task`), routing key `task-created`. The payload is a flat, self-contained JSON fact (`taskId`, `userId`, `description`, `createdAt`) stating what happened and naming no action; consuming services subscribe and decide their own reactions, exactly as this service decides its reaction to invoice-due events. Publication is best-effort after the database commit: a failed publish is logged and swallowed, so task creation never depends on an announcement succeeding.

## Context

The service hears the world but never speaks: it consumes invoice-due events, yet nothing outside it can learn that a task came into existence. Any capability reacting to task activity — digest emails, analytics, audit trails — must poll the HTTP API, read the production database, or be built inside this service, coupling task management to concerns that are not its own. ADR-0004 adopted Watermill deliberately consume-only; this decision reverses that stance with a publisher symmetric to the consumer. Payment reminders are tasks like any other, and the application service seam already erases which transport created a task, so one event covers all creations with no origin marker.

## Considered options

- **Raw amqp091 publishing.** Rejected: reconnection and channel recovery would be hand-rolled again beside a Watermill subscriber that already owns them — two mechanisms for one broker.
- **Watermill AMQP publisher.** Chosen: symmetric with the consumer adoption, over the same broker; messages carry the UUID and JSON content-type headers the consume side already expects.
- **Transactional outbox.** Deferred: an outbox table and relay would buy at-least-once delivery, but with no concrete consumer in existence there is nothing to weigh that guarantee against its cost. Revisited when one exists.
- **Best-effort publish-after-commit.** Chosen: the database write is the truth and the announcement a side effect. Failing task creation — or, on the messaging path, dead-lettering the invoice-due event for redelivery — would create a duplicate task just to retry an announcement for a task that already exists.

## Consequences

- Delivery is at-most-once: a crash between the database commit and the publish loses the event. The transactional outbox is deferred until a concrete consumer exists to demand more.
- A failed publish is logged and swallowed on both transports; the operation's result reflects the database write only. Neither HTTP task creation nor invoice-due handling ever fails or redelivers because of an announcement.
- Replays and duplicates appear honestly in the stream — no dedup, no origin marker; consumers treat every task uniformly and own their reactions.
- The payload is a cross-service contract: the task's database identity, the owner's numeric id (matching the billing contract's user ids, converted at the event boundary), the description (user content, included knowingly as the one field digest and notification consumers need) and the creation time as the service recorded it.
- The `task` exchange is declared by the application's topology declaration alongside the existing exchanges; downstream services own their queues and bindings, and this application declares none.
- Publishing hooks in at the application service seam, so both transports announce creations through one code path; the invoice-due handler calls the same service and carries no publishing logic of its own.
