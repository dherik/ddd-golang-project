package messaging

import (
	"log"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
)

type CalendarQueue struct {
	RabbitMQ rabbitmq.RabbitMQ
}

func NewCalendarQueue(rabbitmq rabbitmq.RabbitMQ) CalendarQueue {
	return CalendarQueue{RabbitMQ: rabbitmq}
}

// func (tq *CalendarQueue) StartListenEvents() error {

// 	conn, err := tq.RabbitMQ.connect()
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	// Create a channel
// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	/*
// 		The queue name start with the name of the service and the event
// 		that the service is listening. This helps to understand the queues that
// 		the service is the owner and which events it's listening just seeing the
// 		list of queues.
// 	*/
// 	q, err := ch.QueueDeclare(
// 		"todo-service.events.calendar.birthday", // name
// 		true,                                    // durable
// 		false,                                   // delete when unused
// 		false,                                   // exclusive
// 		false,                                   // no-wait
// 		nil,                                     // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	/*
// 		Binding the queue with the exchange. The binding is using the routing key "birthday",
// 		what is the kind of calendar event that the todo-service has interest to listen.

// 		The exchange name is the name of the service that is responsible to emit the event. So the
// 		"calendar" exchange is the place where you can see events from the (hypothetical) calendar service.
// 	*/
// 	err = ch.QueueBind(
// 		q.Name,     // queue name
// 		"birthday", // routing key
// 		"calendar", // exchange
// 		false,
// 		nil,
// 	)
// 	failOnError(err, "Failed to bind a queue")

// 	// Consume messages from the queue
// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumerTag
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // arguments
// 	)
// 	failOnError(err, "Failed to register a consumer")

// 	go func() {
// 		for d := range msgs {
// 			// Process the message
// 			slog.Info(fmt.Sprintf("Received message: %s", d.Body))
// 		}
// 	}()

// 	slog.Info("Listening for messages...")

// 	return nil
// }

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
