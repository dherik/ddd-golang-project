package rabbitmq

import (
	"fmt"
	"log"
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
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) DeclareQueueAndBind(queueName, exchange, routingKey string) (amqp.Queue, error) {

	queue, err := r.DeclareQueue(queueName)
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

func (r *RabbitMQ) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
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
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) Close() {
	if err := r.channel.Close(); err != nil {
		log.Panic(fmt.Errorf("failed to close channel: %w", err))
	}
	if err := r.conn.Close(); err != nil {
		log.Panic(fmt.Errorf("failed to close connection: %w", err))
	}
}
