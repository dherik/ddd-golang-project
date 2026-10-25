package messaging

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill"
	wamqp "github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dherik/ddd-golang-project/internal/app/task"
)

// jsonMarshaler publishes the payload with the headers the consume side
// already expects: the message UUID under the toolkit's header key and a JSON
// content type on top of the payloadMarshaler's pass-through body.
var jsonMarshaler = wamqp.DefaultMarshaler{
	PostprocessPublishing: func(publishing amqp.Publishing) amqp.Publishing {
		publishing.ContentType = "application/json"
		return publishing
	},
}

// Publisher owns event publication: a Watermill AMQP publisher over its own
// connection, symmetric with the consumer adoption (see ADR-0004 and
// ADR-0005). The exchange is pinned in the config and the topology is
// application-declared, so the publisher only publishes.
type Publisher struct {
	publisher *wamqp.Publisher
}

func NewPublisher(amqpURI string) (*Publisher, error) {
	logger := watermill.NewSlogLogger(slog.Default())

	config := wamqp.Config{
		Connection: wamqp.ConnectionConfig{AmqpURI: amqpURI},
		Marshaler:  jsonMarshaler,
		Exchange: wamqp.ExchangeConfig{
			// The application declares the task exchange; the publisher
			// announces facts on it without declaring anything.
			GenerateName: wamqp.GenerateExchangeNameConstant(ExchangeTask),
		},
		Publish: wamqp.PublishConfig{
			// Every task created fact routes by the one routing key.
			GenerateRoutingKey: func(string) string { return RoutingKeyTaskCreated },
		},
		TopologyBuilder: existingTopology{},
	}

	publisher, err := wamqp.NewPublisher(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed creating the AMQP publisher: %w", err)
	}

	return &Publisher{publisher: publisher}, nil
}

// PublishTaskCreated publishes one task created fact on the task exchange
// with the task-created routing key.
func (p *Publisher) PublishTaskCreated(event task.TaskCreated) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed encoding the task created event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	if err := p.publisher.Publish(RoutingKeyTaskCreated, msg); err != nil {
		return fmt.Errorf("failed publishing the task created event: %w", err)
	}
	return nil
}

// Close releases the publisher's connection.
func (p *Publisher) Close() error {
	return p.publisher.Close()
}
