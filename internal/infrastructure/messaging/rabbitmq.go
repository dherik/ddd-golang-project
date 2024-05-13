package messaging

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
	RabbitMQDataSource RabbitMQDataSource
}

func NewRabbitMQ(rabbitMQDataSource RabbitMQDataSource) RabbitMQ {
	return RabbitMQ{
		RabbitMQDataSource: rabbitMQDataSource,
	}
}

func (r *RabbitMQ) connect() (*amqp.Connection, error) {

	conn, err := amqp.Dial(r.RabbitMQDataSource.ConnectionString())

	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ : %v", err)
		return &amqp.Connection{}, err
	}
	return conn, nil
}
