package rabbitmq

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQDataSource struct {
	Host     string
	Port     int
	User     string
	Password string
	AmqpURL  string
}

func (ds *RabbitMQDataSource) ConnectionString() string {
	if ds.AmqpURL != "" {
		return ds.AmqpURL
	}
	port := strconv.Itoa(ds.Port)
	return "amqp://" + ds.User + ":" + ds.Password + "@" + ds.Host + ":" + port
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ %s: %w", url, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil {
			slog.Error("failed to close RabbitMQ connection after channel error", "err", closeErr)
		}
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) DeclareQueueAndBind(queueName, exchange, routingKey string, args amqp.Table) (amqp.Queue, error) {

	queue, err := r.DeclareQueue(queueName, args)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare a queue %s: %w", queueName, err)
	}

	err = r.BindQueue(queue.Name, exchange, routingKey)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to bind a queue %s in the exchange %s with routing key %s:  %w",
			queueName, exchange, routingKey, err)
	}

	return queue, nil
}

// DeclareDeadLetterQueue declares the dead-letter exchange, a dead-letter
// queue, and the binding between them. Dead-lettered events are parked
// there for inspection (see docs/adr/0003).
func (r *RabbitMQ) DeclareDeadLetterQueue(exchange, queueName, routingKey string) error {
	err := r.channel.ExchangeDeclare(
		exchange,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare the dead-letter exchange %s: %w", exchange, err)
	}

	_, err = r.DeclareQueue(queueName, nil)
	if err != nil {
		return fmt.Errorf("failed to declare the dead-letter queue %s: %w", queueName, err)
	}

	if err := r.BindQueue(queueName, exchange, routingKey); err != nil {
		return fmt.Errorf("failed to bind the dead-letter queue %s to the exchange %s with routing key %s: %w",
			queueName, exchange, routingKey, err)
	}

	return nil
}

// DeadLetterArgs returns the queue arguments that route events nacked
// without requeue to the given dead-letter exchange. The broker keeps the
// original routing key, so each dead-letter queue binds with its work
// queue's routing key.
func DeadLetterArgs(deadLetterExchange string) amqp.Table {
	return amqp.Table{
		"x-dead-letter-exchange": deadLetterExchange,
	}
}

func (r *RabbitMQ) DeclareQueue(queueName string, args amqp.Table) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)
}

func (r *RabbitMQ) BindQueue(queueName, exchange, routingKey string) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		exchange,
		false,
		nil,
	)
}

func (r *RabbitMQ) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			slog.Error("failed to close RabbitMQ channel", "err", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			slog.Error("failed to close RabbitMQ connection", "err", err)
		}
	}
}
