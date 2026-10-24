# Consume events with Watermill

Status: accepted

Event consumption runs on Watermill: its router executes our handlers, and its AMQP subscriber (watermill-amqp, over the same amqp091-go client we already use) owns the consume loop, reconnection and channel recovery. The retry middleware and our permanent/transient classification land with the follow-up retry ticket; until then every handler failure is nacked without requeue and parks immediately. Adoption is consume-only — no Watermill publisher, no CQRS or event-bus components. Our topology is untouched: the application keeps declaring every exchange, queue and binding, and Watermill consumes on the existing queues under pinned names.

## Context

The hand-written listener carried manual-ack bookkeeping correctly but was about to grow reconnection and bounded retry by hand. Watermill supplies both as battle-tested middleware, and its message abstraction is constructible in tests with injected acknowledgement behaviour. The cost is a dependency and a nomenclature layer (topics mapped onto exchanges, routing keys and queues through generator functions), which is pinned down by configuration and asserted by the integration suite.

## Considered options

- **Keep the hand-written listener.** Rejected: reconnection and retry edge cases become ours to discover in production.
- **wagslane/go-rabbitmq.** Viable and smaller in conceptual footprint (RabbitMQ-native, handler returns an Ack/NackDiscard action), but broker-coupled, single-maintainer, and it would not have carried the middleware ecosystem. Rejected in favour of Watermill's portability.
- **Watermill with its PoisonQueue middleware.** Rejected: PoisonQueue republishes a *copy* of the failure to a topic, changing the dead-letter mechanism ADR-0003 recorded — events must park intact via the broker's dead-letter exchange.
- **Watermill consume-only over our existing topology.** Chosen.

## Consequences

- `watermill` and `watermill-amqp/v3` are added; `amqp091-go` remains a direct dependency because the topology declaration code still uses it.
- Every ADR-0003 decision carries over — manual ack, three-attempt transient retry, poison parked immediately, one dead-letter exchange with per-queue dead-letter queues, application-declared topology. Only the mechanism behind them changes; the retry bound and classification arrive with the follow-up ticket.
- Reconnection and channel recovery are watermill-amqp's responsibility; the hand-rolled reconnect work is retired.
- Handlers no longer own acknowledgement — the router acks on success and nacks on error — so the hand-rolled delivery abstraction is deleted and test messages are built in one line with no fake.
- A payload pass-through marshaler replaces the toolkit's default: only the payload is our contract, so a non-string header must never turn a valid event into poison.
- Watermill's topic-to-queue name mapping is pinned by generators; a mismatch fails against the durable queues rather than silently redeclaring them, and the integration suite exercises the exact queue and dead-letter names end to end.
- Watermill's logger adapter wraps our slog logger so consumption logs stay in one format. The router is cancellation-aware: cancelling its context stops consumption cleanly; wiring process signals to that cancellation remains a separate lifecycle concern.

Update: the consume-only scope was superseded for publishing by [ADR-0005](0005-publish-task-created-events.md) — a Watermill AMQP publisher now announces task-created events over the same broker. The narrower stance still holds for the rest: no CQRS or event-bus components.
