# Dead-letter failed event consumption

Status: accepted

Event consumers run manual acknowledgement. A handler success is acked; a permanent failure (poison message) is nacked straight to a dead-letter queue; a transient failure is retried a bounded number of times in-process, then dead-lettered. One dead-letter exchange serves all work queues, each with its own dead-letter queue, all declared by the application at startup. Each listener owns its channel and reconnects with backoff instead of exiting silently.

## Context

The invoice-due consumer (ADR-0002) shipped auto-ack with log-and-drop: a handler failure, or a crash mid-processing, silently loses the message. Reviewing the listener also showed both consumers multiplexing a single AMQP channel and exiting without a sound when it closes. Once real billing events flow, silent loss is not an acceptable failure mode — an invoice reminder that vanishes is indistinguishable from one that never existed.

## Considered options

- **Keep auto-ack and drop on failure.** Rejected: silent loss, and nothing to inspect afterwards.
- **TTL retry-queue topology** (retry queue dead-lettering back to the work queue after a delay). Rejected for now: doubles the queue count and needs delay tuning before any retry evidence exists. In-process retry covers the expected blip.
- **A dead-letter exchange per work queue.** Rejected: one direct exchange with per-queue routing keys gives the same isolation with less topology.

## Consequences

- Nothing vanishes silently: poison events are parked in a dead-letter queue for inspection. Someone must look — there is no consumer, replay tooling, or alerting on dead-letter queues in v1.
- Manual ack means unacknowledged messages are redelivered after a restart, so a crash mid-processing no longer loses an event.
- A transient outage delays a queue by at most the retry window; a poison message is dead-lettered immediately without burning retries.
- One channel per listener: a channel error no longer silences the other consumer, and the listener reconnects instead of dying.
- This supersedes two consequences recorded in ADR-0002: consumers are no longer auto-ack, and invalid invoice-due events are no longer "logged and skipped" — they are dead-lettered. ADR-0002's core decision (react to billing events by creating reminder tasks) is unchanged.
