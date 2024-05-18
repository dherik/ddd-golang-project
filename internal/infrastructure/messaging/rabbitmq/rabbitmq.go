package rabbitmq

import (
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
	// RabbitMQDataSource RabbitMQDataSource
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// return RabbitMQ{
	// 	RabbitMQDataSource: rabbitMQDataSource,
	// }

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

// func (r *RabbitMQ) connect() (*amqp.Connection, error) {

// 	conn, err := amqp.Dial(r.RabbitMQDataSource.ConnectionString())

// 	if err != nil {
// 		log.Fatalf("failed to connect to RabbitMQ : %v", err)
// 		return &amqp.Connection{}, err
// 	}
// 	return conn, nil
// }

func (r *RabbitMQ) DeclareQueueAndBind(queueName, exchange, routingKey string) (amqp.Queue, error) {

	queue, err := r.DeclareQueue(queueName)
	failOnError(err, "Failed to declare a queue") //TODO log queue name

	err = r.BindQueue(queue.Name, exchange, routingKey)
	failOnError(err, "Failed to bind a queue") //TODO log binding details

	return queue, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
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
		log.Println("Error closing channel:", err) //FIXME
	}
	if err := r.conn.Close(); err != nil {
		log.Println("Error closing connection:", err) //FIXME
	}
}
