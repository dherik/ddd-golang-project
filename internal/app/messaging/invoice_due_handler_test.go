package messaging

import (
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDueEventCreatesReminderTask(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
		return task.UserId == "1" && task.Description == "Pay invoice INV-42 by 2026-11-01"
	})).Return(domain.Task{}, nil)

	handler := NewInvoiceDueHandler(api.NewTaskService(mockRepo))

	body := []byte(`{"invoiceId":"INV-42","userId":1,"dueDate":"2026-11-01"}`)
	err := handler.HandleMessage(context.Background(), amqp.Delivery{Body: body})

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestInvoiceDueEventReminderWording(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		description string
	}{
		{
			name:        "issuer only",
			payload:     `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","issuer":"Acme"}`,
			description: "Pay invoice INV-1 (issuer: Acme) by 2026-11-01",
		},
		{
			name:        "amount and currency only",
			payload:     `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","amount":"199.90","currency":"EUR"}`,
			description: "Pay invoice INV-1 (amount: 199.90 EUR) by 2026-11-01",
		},
		{
			name:        "all optional fields",
			payload:     `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","issuer":"Acme","amount":"199.90","currency":"EUR"}`,
			description: "Pay invoice INV-1 (issuer: Acme, amount: 199.90 EUR) by 2026-11-01",
		},
		{
			name:        "amount without currency is omitted",
			payload:     `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","amount":"199.90"}`,
			description: "Pay invoice INV-1 by 2026-11-01",
		},
		{
			name:        "currency without amount is omitted",
			payload:     `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","currency":"EUR"}`,
			description: "Pay invoice INV-1 by 2026-11-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockTaskRepository)
			mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
				return task.Description == tt.description
			})).Return(domain.Task{}, nil)

			handler := NewInvoiceDueHandler(api.NewTaskService(mockRepo))

			err := handler.HandleMessage(context.Background(), amqp.Delivery{Body: []byte(tt.payload)})

			require.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestInvoiceDueEventInvalidPayloadIsSkipped(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{name: "not JSON", payload: `this is not json`},
		{name: "empty invoice id", payload: `{"invoiceId":"","userId":1,"dueDate":"2026-11-01"}`},
		{name: "missing invoice id", payload: `{"userId":1,"dueDate":"2026-11-01"}`},
		{name: "zero user id", payload: `{"invoiceId":"INV-1","userId":0,"dueDate":"2026-11-01"}`},
		{name: "negative user id", payload: `{"invoiceId":"INV-1","userId":-3,"dueDate":"2026-11-01"}`},
		{name: "missing user id", payload: `{"invoiceId":"INV-1","dueDate":"2026-11-01"}`},
		{name: "unparseable due date", payload: `{"invoiceId":"INV-1","userId":1,"dueDate":"next tuesday"}`},
		{name: "missing due date", payload: `{"invoiceId":"INV-1","userId":1}`},
		{name: "datetime instead of date", payload: `{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01T10:00:00Z"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockTaskRepository)
			handler := NewInvoiceDueHandler(api.NewTaskService(mockRepo))

			err := handler.HandleMessage(context.Background(), amqp.Delivery{Body: []byte(tt.payload)})

			require.NoError(t, err)
			mockRepo.AssertNotCalled(t, "AddTaskToUser", mock.Anything, mock.Anything)
		})
	}
}
