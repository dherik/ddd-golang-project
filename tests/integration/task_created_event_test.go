package integration

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	wamqp "github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/internal/app/messaging"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
)

// taskCreatedTestQueue is the test-owned queue the task created contract is
// asserted on. The application declares no downstream queues (ADR-0005); this
// one exists only for this test.
const taskCreatedTestQueue = "todo-service.test.task-created"

// The contract test for ADR-0005: a task created through the real HTTP API
// arrives on the task exchange as exactly one flat JSON fact, on a queue this
// test owns.
func (s *TaskTestSuite) TestTaskCreatedEventPublishedWhenTaskCreatedOverHTTP() {
	s.Require().NoError(setup.BindTestQueue(taskCreatedTestQueue,
		messaging.ExchangeTask, messaging.RoutingKeyTaskCreated))

	token := setup.Login("admin", "some_password")
	description := "Announced as a task created fact"

	response := addTask(s, api.TaskRequest{Description: description}, token)
	s.Equal(http.StatusCreated, response.StatusCode)
	response.Body.Close()

	created := findTaskByDescription(s, token, description)
	s.Require().NotNil(created, "the created task should be listed for its owner")

	event := waitForTaskCreatedEvent(s, description)

	// The payload is the flat fact: taskId, userId, description, createdAt —
	// nothing else, no instruction, no origin marker.
	var fact struct {
		TaskId      int       `json:"taskId"`
		UserId      int       `json:"userId"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"createdAt"`
	}
	s.NoError(json.Unmarshal([]byte(event.Payload), &fact))
	s.Equal(created.Id, fact.TaskId)
	s.Equal(1, fact.UserId, "the owner id is numeric, matching the billing contract")
	s.Equal(description, fact.Description)
	s.True(fact.CreatedAt.Equal(created.CreatedAt),
		"event createdAt %v should be the task's recorded creation time %v", fact.CreatedAt, created.CreatedAt)

	var raw map[string]any
	s.NoError(json.Unmarshal([]byte(event.Payload), &raw))
	s.ElementsMatch([]string{"taskId", "userId", "description", "createdAt"},
		slices.Collect(maps.Keys(raw)))

	// The creation time is an RFC 3339 timestamp in UTC.
	rawCreatedAt, ok := raw["createdAt"].(string)
	s.Require().True(ok, "createdAt %v should be a JSON string", raw["createdAt"])
	s.True(strings.HasSuffix(rawCreatedAt, "Z"), "createdAt %q should be UTC", rawCreatedAt)
	_, err := time.Parse(time.RFC3339, rawCreatedAt)
	s.NoError(err)

	// Headers: the message UUID the consume side expects, and a JSON content
	// type.
	msgUUID, ok := event.Headers[wamqp.DefaultMessageUUIDHeaderKey].(string)
	s.Require().True(ok, "message uuid header %v should be a string",
		event.Headers[wamqp.DefaultMessageUUIDHeaderKey])
	s.NotEmpty(msgUUID)
	s.Equal("application/json", event.ContentType)
}

// waitForTaskCreatedEvent polls the test queue until exactly one task created
// fact names the description, or fails the test after a deadline. Publishing
// and broker routing are asynchronous relative to the HTTP response, so a
// single immediate read would race the announcement.
func waitForTaskCreatedEvent(s *TaskTestSuite, description string) setup.BrokerMessage {
	s.T().Helper()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		messages, err := setup.ReadTestQueue(taskCreatedTestQueue)
		s.NoError(err)

		var matches []setup.BrokerMessage
		for _, m := range messages {
			if strings.Contains(m.Payload, description) {
				matches = append(matches, m)
			}
		}
		if len(matches) == 1 {
			return matches[0]
		}
		if len(matches) > 1 {
			s.T().Fatalf("expected exactly one task created fact for %q, got %d", description, len(matches))
		}
		time.Sleep(100 * time.Millisecond)
	}

	s.T().Fatalf("no task created fact for %q arrived within the deadline", description)
	return setup.BrokerMessage{}
}
