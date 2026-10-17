package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
)

func (s *TaskTestSuite) TestInvoiceDueEventCreatesPaymentReminder() {
	token := setup.Login("admin", "some_password")

	err := setup.PublishMessage("billing", "invoice-due",
		`{"invoiceId":"INV-42","userId":1,"dueDate":"2026-11-01"}`)
	s.NoError(err)

	task := waitForPaymentReminder(s, token, "Pay invoice INV-42 by 2026-11-01")
	s.Equal("1", task.UserId)
}

func (s *TaskTestSuite) TestInvoiceDueEventReminderIncludesIssuerAndAmount() {
	token := setup.Login("admin", "some_password")

	err := setup.PublishMessage("billing", "invoice-due",
		`{"invoiceId":"INV-7","userId":1,"dueDate":"2026-11-01","issuer":"Acme","amount":"199.90","currency":"EUR"}`)
	s.NoError(err)

	waitForPaymentReminder(s, token, "Pay invoice INV-7 (issuer: Acme, amount: 199.90 EUR) by 2026-11-01")
}

func (s *TaskTestSuite) TestInvoiceDueEventReminderOnlyForNamedUser() {
	err := setup.PublishMessage("billing", "invoice-due",
		`{"invoiceId":"INV-9","userId":1,"dueDate":"2026-11-01"}`)
	s.NoError(err)

	waitForPaymentReminder(s, setup.Login("admin", "some_password"),
		"Pay invoice INV-9 by 2026-11-01")

	otherToken := setup.Login("john_doe", "some_password")
	s.Nil(findTaskByDescription(s, otherToken, "Pay invoice INV-9 by 2026-11-01"))
}

// waitForPaymentReminder polls the task list until a task with the given
// description appears, or fails the test after a deadline. Consumption is
// asynchronous, so a single immediate read would race the listener.
func waitForPaymentReminder(s *TaskTestSuite, token, description string) *api.TaskResponse {
	s.T().Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if task := findTaskByDescription(s, token, description); task != nil {
			return task
		}
		time.Sleep(100 * time.Millisecond)
	}

	s.T().Fatalf("payment reminder %q did not appear within the deadline", description)
	return nil
}

func findTaskByDescription(s *TaskTestSuite, token, description string) *api.TaskResponse {
	s.T().Helper()

	url := fmt.Sprintf("http://localhost:%d/tasks", 3333)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	var page api.TasksResponse
	s.NoError(json.Unmarshal(body, &page))

	for i := range page.Items {
		if page.Items[i].Description == description {
			return &page.Items[i]
		}
	}
	return nil
}
