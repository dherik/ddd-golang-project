package api

import (
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTaskRepository struct {
	mock.Mock
}

func (m *mockTaskRepository) GetByID(id int) (domain.Task, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskRepository) GetByUserID(userId string) ([]domain.Task, error) {
	args := m.Called(userId)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *mockTaskRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	args := m.Called(userId, task)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func TestTaskService_FindTasks(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	expectedTasks := []domain.Task{
		{Id: 1, UserId: "1", Description: "Task 1", CreatedAt: time.Now().UTC()},
		{Id: 2, UserId: "2", Description: "Task 2", CreatedAt: time.Now().UTC()},
	}

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", start, end).Return(expectedTasks, nil)

	svc := NewTaskService(mockRepo)
	tasks, err := svc.FindTasks(start, end)

	require.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, 1, tasks[0].Id)
	assert.Equal(t, "Task 1", tasks[0].Description)
	assert.Equal(t, "1", tasks[0].UserId)
	mockRepo.AssertExpectations(t)
}

func TestTaskService_FindTasks_RepositoryError(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", start, end).Return([]domain.Task{}, assert.AnError)

	svc := NewTaskService(mockRepo)
	_, err := svc.FindTasks(start, end)

	require.Error(t, err)
	mockRepo.AssertExpectations(t)
}
