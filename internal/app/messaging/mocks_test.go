package messaging

import (
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/task"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/mock"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type mockTaskRepository struct {
	mock.Mock
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

type mockTaskPublisher struct {
	mock.Mock
}

func (m *mockTaskPublisher) PublishTaskCreated(event task.TaskCreated) error {
	args := m.Called(event)
	return args.Error(0)
}

// silentTaskPublisher accepts any announcement and succeeds, for tests whose
// subject is the reminder task rather than the announcement itself.
func silentTaskPublisher() *mockTaskPublisher {
	m := new(mockTaskPublisher)
	m.On("PublishTaskCreated", mock.Anything).Return(nil)
	return m
}
