package messaging

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/dherik/ddd-golang-project/internal/app/task"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDueEventCreatesReminderTask(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
		return task.UserId == "1" && task.Description == "Pay invoice INV-42 by 2026-11-01"
	})).Return(domain.Task{}, nil)

	handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, silentTaskPublisher()))

	err := handler.HandleMessage(testMessage(`{"invoiceId":"INV-42","userId":1,"dueDate":"2026-11-01"}`))

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// One handled invoice due event announces its reminder task exactly once: a
// second announcement on a single event would be a duplicate fact for
// consumers counting on the task exchange.
func TestInvoiceDueEventAnnouncesReminderTaskOnce(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.Anything).Return(domain.Task{
		Id:          7,
		UserId:      "1",
		Description: "Pay invoice INV-42 by 2026-11-01",
	}, nil)

	publisher := new(mockTaskPublisher)
	publisher.On("PublishTaskCreated", mock.Anything).Return(nil).Once()

	handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, publisher))

	err := handler.HandleMessage(testMessage(`{"invoiceId":"INV-42","userId":1,"dueDate":"2026-11-01"}`))

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	publisher.AssertExpectations(t)
}

// A failed publish must not fail the handler: failing it would nack the
// invoice due event into the dead-letter queue, and its redelivery would
// create a duplicate reminder just to retry an announcement for a task that
// already exists.
func TestInvoiceDueEventReminderCreatedWhenPublishingFails(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
		return task.UserId == "1" && task.Description == "Pay invoice INV-42 by 2026-11-01"
	})).Return(domain.Task{
		Id:          7,
		UserId:      "1",
		Description: "Pay invoice INV-42 by 2026-11-01",
	}, nil)

	failingPublisher := new(mockTaskPublisher)
	failingPublisher.On("PublishTaskCreated", mock.Anything).Return(errors.New("broker unavailable"))

	handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, failingPublisher))

	err := handler.HandleMessage(testMessage(`{"invoiceId":"INV-42","userId":1,"dueDate":"2026-11-01"}`))

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	failingPublisher.AssertExpectations(t)
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

			handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, silentTaskPublisher()))

			err := handler.HandleMessage(testMessage(tt.payload))

			require.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestInvoiceDueEventInvalidPayloadIsPermanent(t *testing.T) {
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
			handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, silentTaskPublisher()))

			err := handler.HandleMessage(testMessage(tt.payload))

			require.Error(t, err)
			require.True(t, IsPermanent(err), "expected a permanent failure, got: %v", err)
			mockRepo.AssertNotCalled(t, "AddTaskToUser", mock.Anything, mock.Anything)
		})
	}
}

func TestInvoiceDueEventConstructionAndPersistenceFailuresAreTransient(t *testing.T) {
	t.Run("task construction failure", func(t *testing.T) {
		// A 2000-character issuer pushes the rendered description past the
		// domain's 1024-character limit, failing construction.
		longIssuer := strings.Repeat("A", 2000)
		payload := fmt.Sprintf(`{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01","issuer":%q}`, longIssuer)

		mockRepo := new(mockTaskRepository)
		handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, silentTaskPublisher()))

		err := handler.HandleMessage(testMessage(payload))

		require.Error(t, err)
		require.False(t, IsPermanent(err), "construction failures stay transient")
	})

	t.Run("persistence failure", func(t *testing.T) {
		mockRepo := new(mockTaskRepository)
		mockRepo.On("AddTaskToUser", mock.Anything, mock.Anything).
			Return(domain.Task{}, fmt.Errorf("database unavailable"))
		handler := NewInvoiceDueHandler(task.NewService(mockRepo, fixedClock{}, silentTaskPublisher()))

		err := handler.HandleMessage(testMessage(`{"invoiceId":"INV-1","userId":1,"dueDate":"2026-11-01"}`))

		require.Error(t, err)
		require.False(t, IsPermanent(err), "persistence failures stay transient")
	})
}

func testMessage(payload string) *message.Message {
	return message.NewMessage("test-message-id", []byte(payload))
}
