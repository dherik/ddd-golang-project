package messaging

import (
	"time"

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
