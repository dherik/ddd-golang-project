package messaging

import (
	"context"
	"log"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, msg amqp.Delivery) error
}

type MessageListener struct {
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewMessageListener(rabbitMQ *rabbitmq.RabbitMQ) *MessageListener {
	return &MessageListener{rabbitMQ: rabbitMQ}
}

func (ml *MessageListener) StartListening(ctx context.Context, queueName string, handler MessageHandler) error {
	msgs, err := ml.rabbitMQ.Consume(queueName)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				if err := handler.HandleMessage(ctx, msg); err != nil {
					log.Println("Error handling message:", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
