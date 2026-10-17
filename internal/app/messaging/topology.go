package messaging

import (
	"fmt"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
)

// Broker topology names. These are the service's contract with the external
// calendar and billing services and with its own dead-letter queues; every
// other place in the codebase imports them instead of restating the strings.
const (
	ExchangeCalendar     = "calendar"
	ExchangeBilling      = "billing"
	RoutingKeyBirthday   = "birthday"
	RoutingKeyInvoiceDue = "invoice-due"

	QueueCalendarBirthday  = "todo-service.events.calendar.birthday"
	QueueBillingInvoiceDue = "todo-service.events.billing.invoice-due"

	DeadLetterExchange      = "todo-service.dead-letter"
	DeadLetterQueueCalendar = "todo-service.dead.calendar.birthday"
	DeadLetterQueueBilling  = "todo-service.dead.billing.invoice-due"
)

// Declare builds the complete broker topology in dependency order: the
// dead-letter exchange and queues first, then the work exchanges, then the
// work queues bound with dead-letter arguments. The application is the
// single declarer of its topology (ADR-0004); declaration failure is fatal
// at the call site.
func Declare(conn *rabbitmq.RabbitMQ) error {
	if err := conn.DeclareDeadLetterQueue(DeadLetterExchange,
		DeadLetterQueueBilling, RoutingKeyInvoiceDue); err != nil {
		return fmt.Errorf("failed declaring the billing dead-letter queue: %w", err)
	}
	if err := conn.DeclareDeadLetterQueue(DeadLetterExchange,
		DeadLetterQueueCalendar, RoutingKeyBirthday); err != nil {
		return fmt.Errorf("failed declaring the calendar dead-letter queue: %w", err)
	}

	// Direct exchanges: external publishers route by exact routing key
	// (see ADR-0002), so nothing here may drift to another exchange type.
	if err := conn.DeclareExchange(ExchangeCalendar, rabbitmq.ExchangeDirect); err != nil {
		return fmt.Errorf("failed declaring the %s exchange: %w", ExchangeCalendar, err)
	}
	if err := conn.DeclareExchange(ExchangeBilling, rabbitmq.ExchangeDirect); err != nil {
		return fmt.Errorf("failed declaring the %s exchange: %w", ExchangeBilling, err)
	}

	deadLetterArgs := rabbitmq.DeadLetterArgs(DeadLetterExchange)

	if _, err := conn.DeclareQueueAndBind(QueueCalendarBirthday,
		ExchangeCalendar, RoutingKeyBirthday, deadLetterArgs); err != nil {
		return fmt.Errorf("failed declaring the birthday work queue: %w", err)
	}
	if _, err := conn.DeclareQueueAndBind(QueueBillingInvoiceDue,
		ExchangeBilling, RoutingKeyInvoiceDue, deadLetterArgs); err != nil {
		return fmt.Errorf("failed declaring the invoice due work queue: %w", err)
	}

	return nil
}
