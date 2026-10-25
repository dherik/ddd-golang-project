package setup

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// BrokerMessage is one message read off a test-owned queue: the payload plus
// the delivery metadata the published contract is asserted against.
type BrokerMessage struct {
	Payload     string
	ContentType string
	Headers     amqp.Table
}

// BindTestQueue declares a test-owned queue and binds it to the exchange with
// the routing key, so a test can observe what the application publishes. The
// application declares no downstream queues or bindings; a queue bound before
// a publish is what makes that publish observable.
func BindTestQueue(queueName, exchange, routingKey string) error {
	conn, err := amqp.Dial(RabbitMQDataSource.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed connecting to the broker: %w", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed opening a channel: %w", err)
	}

	if _, err := channel.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed declaring the queue %s: %w", queueName, err)
	}
	if err := channel.QueueBind(queueName, routingKey, exchange, false, nil); err != nil {
		return fmt.Errorf("failed binding the queue %s to the exchange %s with routing key %s: %w",
			queueName, exchange, routingKey, err)
	}
	return nil
}

// ReadTestQueue drains a test-owned queue, reading each message exactly once
// (basic.get with auto-ack). Polling callers re-read between publishes, so a
// message is observed on the poll after it arrives at the latest.
func ReadTestQueue(queueName string) ([]BrokerMessage, error) {
	conn, err := amqp.Dial(RabbitMQDataSource.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed connecting to the broker: %w", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed opening a channel: %w", err)
	}

	var messages []BrokerMessage
	for {
		delivery, ok, err := channel.Get(queueName, true)
		if err != nil {
			return nil, fmt.Errorf("failed reading from the queue %s: %w", queueName, err)
		}
		if !ok {
			return messages, nil
		}
		messages = append(messages, BrokerMessage{
			Payload:     string(delivery.Body),
			ContentType: delivery.ContentType,
			Headers:     delivery.Headers,
		})
	}
}
