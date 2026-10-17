package messaging

import (
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/mock"
)

type mockTaskRepository struct {
	mock.Mock
}

// fakeDelivery is a recording stand-in for a broker delivery, satisfying the
// listener's acknowledgement interface without a broker.
type fakeDelivery struct {
	body        string
	acked       bool
	nacked      bool
	nackRequeue bool
}

func (f fakeDelivery) Body() []byte { return []byte(f.body) }

func (f *fakeDelivery) Ack(multiple bool) error {
	f.acked = true
	return nil
}

func (f *fakeDelivery) Nack(multiple, requeue bool) error {
	f.nacked = true
	f.nackRequeue = requeue
	return nil
}

func testDelivery(body string) *fakeDelivery {
	return &fakeDelivery{body: body}
}

func (m *mockTaskRepository) GetByIDAndUserID(id int, userId string) (domain.Task, error) {
	args := m.Called(id, userId)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	args := m.Called(userId, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskRepository) FindTasks(query domain.TaskQuery) (domain.TaskPage, error) {
	args := m.Called(query)
	return args.Get(0).(domain.TaskPage), args.Error(1)
}
