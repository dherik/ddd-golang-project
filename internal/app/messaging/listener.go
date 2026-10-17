package messaging

import (
	"context"
	"log"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
)

type MessageListener struct {
	rabbitMQ *rabbitmq.RabbitMQ
}

func NewMessageListener(rabbitMQ *rabbitmq.RabbitMQ) *MessageListener {
	return &MessageListener{rabbitMQ: rabbitMQ}
}

// brokerDelivery adapts the broker's delivery to the Delivery interface.
type brokerDelivery struct {
	d amqp.Delivery
}

func (b brokerDelivery) Body() []byte { return b.d.Body }

func (b brokerDelivery) Ack(multiple bool) error { return b.d.Ack(multiple) }

func (b brokerDelivery) Nack(multiple, requeue bool) error { return b.d.Nack(multiple, requeue) }

// StartListening subscribes to a queue and consumes it with manual
// acknowledgement: the listener owns every ack and nack decision.
func (ml *MessageListener) StartListening(ctx context.Context, queueName string, handler MessageHandler) error {
	msgs, err := ml.rabbitMQ.Consume(queueName)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ message channel closed")
					return
				}
				processDelivery(ctx, handler, brokerDelivery{d: msg})
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// processDelivery applies the consumption policy to one event: acknowledge
// success, dead-letter a permanent failure immediately, and dead-letter a
// transient failure (in-process retries arrive with the retry ticket).
func processDelivery(ctx context.Context, handler MessageHandler, msg Delivery) {
	err := handler.HandleMessage(ctx, msg)
	switch {
	case err == nil:
		if ackErr := msg.Ack(false); ackErr != nil {
			slog.Error("failed acknowledging event", "err", ackErr)
		}
	case IsPermanent(err):
		slog.Warn("dead-lettering event after permanent failure", "err", err)
		if nackErr := msg.Nack(false, false); nackErr != nil {
			slog.Error("failed dead-lettering event", "err", nackErr)
		}
	default:
		slog.Error("dead-lettering event after transient failure", "err", err)
		if nackErr := msg.Nack(false, false); nackErr != nil {
			slog.Error("failed dead-lettering event", "err", nackErr)
		}
	}
}
