package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQDataSource struct {
	Host     string
	Port     int
	User     string
	Password string
}

type RabbitMQ struct {
}

func NewRabbitMQ() RabbitMQ {
	return RabbitMQ{}
}

func (r *RabbitMQ) connect() (*amqp.Connection, error) {

	conn, err := amqp.Dial(
		"amqp://guest:guest@localhost", //FIXME
	)
	if err != nil {
		// log.Fatal("deu ruim")
		log.Fatal(err) //FIXME
		return &amqp.Connection{}, err
	}
	// defer conn.Close()

	return conn, nil

	// consumer, err := rabbitmq.NewConsumer(
	// 	conn,
	// 	"my_queue",
	// 	rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
	// 	rabbitmq.WithConsumerOptionsExchangeName("events"),
	// 	rabbitmq.WithConsumerOptionsExchangeDeclare,
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer consumer.Close()

	// err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
	// 	log.Printf("consumed: %v", string(d.Body))
	// 	// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
	// 	return rabbitmq.Ack
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
