package messaging

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CalendarHandler struct{}

func (h *CalendarHandler) HandleMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Println("Received calendar message:", string(msg.Body))
	// Handle the message (e.g., parse JSON, update calendar, etc.)
	return nil
}
