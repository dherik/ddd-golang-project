package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/internal/domain"
)

const dateOnlyLayout = "2006-01-02"

// InvoiceDueEvent is the payload published by the external billing service
// when an invoice belonging to a user must be paid by a due date. It is a
// fact, not an instruction: what this service does with it is its own
// decision (see docs/adr/0002).
type InvoiceDueEvent struct {
	InvoiceID string `json:"invoiceId"`
	UserID    int    `json:"userId"`
	DueDate   string `json:"dueDate"`
	Issuer    string `json:"issuer"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

type InvoiceDueHandler struct {
	TaskService *api.TaskService
}

func NewInvoiceDueHandler(taskService *api.TaskService) *InvoiceDueHandler {
	return &InvoiceDueHandler{TaskService: taskService}
}

func (h *InvoiceDueHandler) HandleMessage(ctx context.Context, msg amqp.Delivery) error {
	var event InvoiceDueEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		slog.Warn("discarding invoice due event with invalid payload", "err", err)
		return nil
	}

	if err := event.Validate(); err != nil {
		slog.Warn("discarding invoice due event with invalid fields",
			"err", err, "invoiceId", event.InvoiceID, "userId", event.UserID, "dueDate", event.DueDate)
		return nil
	}

	task, err := domain.NewTask(strconv.Itoa(event.UserID), event.ReminderDescription())
	if err != nil {
		slog.Error("failed building payment reminder task", "err", err)
		return nil
	}

	if err := h.TaskService.AddTaskToUser(task); err != nil {
		slog.Error("failed persisting payment reminder task", "err", err)
		return nil
	}

	slog.Info("created payment reminder task",
		"invoiceId", event.InvoiceID, "userId", event.UserID, "dueDate", event.DueDate)
	return nil
}

// Validate checks the fields the billing service must always send. The due
// date is a calendar date; a timestamp is rejected rather than truncated.
func (e InvoiceDueEvent) Validate() error {
	if e.InvoiceID == "" {
		return errors.New("invoiceId is required")
	}
	if e.UserID <= 0 {
		return errors.New("userId must be a positive integer")
	}
	if _, err := time.Parse(dateOnlyLayout, e.DueDate); err != nil {
		return fmt.Errorf("dueDate must be a calendar date (%s): %w", dateOnlyLayout, err)
	}
	return nil
}

// ReminderDescription renders the payment reminder wording: the invoice and
// due date always, plus a parenthetical with the details the billing service
// knows. Amount appears only together with its currency.
func (e InvoiceDueEvent) ReminderDescription() string {
	var parts []string
	if e.Issuer != "" {
		parts = append(parts, fmt.Sprintf("issuer: %s", e.Issuer))
	}
	if e.Amount != "" && e.Currency != "" {
		parts = append(parts, fmt.Sprintf("amount: %s %s", e.Amount, e.Currency))
	}
	details := ""
	if len(parts) > 0 {
		details = fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}
	return fmt.Sprintf("Pay invoice %s%s by %s", e.InvoiceID, details, e.DueDate)
}

type CalendarHandler struct{}

func (h *CalendarHandler) HandleMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Println("Received calendar message:", string(msg.Body))
	// Handle the message (e.g., parse JSON, update calendar, etc.)
	return nil
}
