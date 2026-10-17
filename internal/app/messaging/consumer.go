package messaging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	wamqp "github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	amqp "github.com/rabbitmq/amqp091-go"
)

// payloadMarshaler passes the event payload through untouched. Only the
// payload is part of our contract with the billing service — headers are
// not — so a non-string header must never turn a valid event into poison,
// which the toolkit's DefaultMarshaler would do.
type payloadMarshaler struct{}

func (payloadMarshaler) Marshal(msg *message.Message) (amqp.Publishing, error) {
	return amqp.Publishing{Body: msg.Payload, DeliveryMode: amqp.Persistent}, nil
}

func (payloadMarshaler) Unmarshal(delivery amqp.Delivery) (*message.Message, error) {
	return message.NewMessage(watermill.NewUUID(), delivery.Body), nil
}

// existingTopology leaves topology declaration to the application: the
// toolkit consumes on queues the service has already declared and never
// declares, binds, or renames anything itself.
type existingTopology struct{}

func (existingTopology) BuildTopology(
	channel *amqp.Channel,
	params wamqp.BuildTopologyParams,
	config wamqp.Config,
	logger watermill.LoggerAdapter,
) error {
	return nil
}

func (existingTopology) ExchangeDeclare(channel *amqp.Channel, exchangeName string, config wamqp.Config) error {
	return nil
}

// Consumer owns event consumption: a Watermill router runs the handlers, and
// the AMQP subscriber owns the consume loop, reconnection and channel
// recovery. Queue names are pinned to the topic and the topology is
// application-declared, so the subscriber only consumes (see ADR-0003 and
// ADR-0004).
type Consumer struct {
	router     *message.Router
	subscriber *wamqp.Subscriber
}

func NewConsumer(amqpURI string) (*Consumer, error) {
	logger := watermill.NewSlogLogger(slog.Default())

	queueName := wamqp.GenerateQueueNameTopicName
	identityRoutingKey := func(topic string) string { return topic }

	config := wamqp.Config{
		Connection: wamqp.ConnectionConfig{AmqpURI: amqpURI},
		Marshaler:  payloadMarshaler{},
		Exchange: wamqp.ExchangeConfig{
			// No exchange: the application binds the queue already.
			GenerateName: wamqp.GenerateExchangeNameConstant(""),
		},
		Queue: wamqp.QueueConfig{
			// The topic is the queue name, so consumption targets exactly
			// the queue the application declared.
			GenerateName: queueName,
		},
		QueueBind: wamqp.QueueBindConfig{
			GenerateRoutingKey: identityRoutingKey,
		},
		Consume: wamqp.ConsumeConfig{
			// Nack without requeue: the broker routes failures to the
			// queue's dead-letter exchange, where they park intact.
			NoRequeueOnNack: true,
			Qos:             wamqp.QosConfig{PrefetchCount: 1},
		},
		TopologyBuilder: existingTopology{},
	}

	subscriber, err := wamqp.NewSubscriber(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed creating the AMQP subscriber: %w", err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed creating the event router: %w", err)
	}

	return &Consumer{router: router, subscriber: subscriber}, nil
}

// Add registers a handler for a queue. The topic is the queue name, so the
// pinned generator consumes exactly that queue.
func (c *Consumer) Add(queueName string, handler func(msg *message.Message) error) {
	c.router.AddConsumerHandler(queueName, queueName, c.subscriber, handler)
}

// Run consumes until the context is cancelled. It blocks; run it in a
// goroutine. Cancellation stops consumption gracefully, leaving an in-flight
// event finished or unacknowledged for redelivery.
func (c *Consumer) Run(ctx context.Context) error {
	return c.router.Run(ctx)
}
