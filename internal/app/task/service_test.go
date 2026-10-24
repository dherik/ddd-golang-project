package task

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

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

var testClock = fixedClock{now: time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)}

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

func TestService_CreateTask(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
		return task.UserId == "1" &&
			task.Description == "some description" &&
			task.CreatedAt.Equal(testClock.now)
	// the persisted task carries the database-generated id; the constructed one does not
	})).Return(domain.Task{
		Id:          7,
		UserId:      "1",
		Description: "some description",
		CreatedAt:   testClock.now,
	}, nil)

	err := NewService(mockRepo, testClock).CreateTask("1", "some description")

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_CreateTask_InvalidDescription(t *testing.T) {
	mockRepo := new(mockTaskRepository)

	err := NewService(mockRepo, testClock).CreateTask("1", "")

	require.ErrorIs(t, err, domain.ErrDescriptionInvalid)
	mockRepo.AssertNotCalled(t, "AddTaskToUser", mock.Anything, mock.Anything)
}

func TestService_FindTasks(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	query := domain.TaskQuery{UserId: "1", StartDate: &start, EndDate: &end, PageSize: 20}

	expectedTasks := []domain.Task{
		{Id: 1, UserId: "1", Description: "Task 1", CreatedAt: time.Now().UTC()},
	}

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", query).Return(domain.TaskPage{Tasks: expectedTasks, HasMore: false}, nil)

	svc := NewService(mockRepo, testClock)
	page, err := svc.FindTasks(query)

	require.NoError(t, err)
	assert.Len(t, page.Tasks, 1)
	assert.False(t, page.HasMore)
	assert.Equal(t, 1, page.Tasks[0].Id)
	assert.Equal(t, "Task 1", page.Tasks[0].Description)
	assert.Equal(t, "1", page.Tasks[0].UserId)
	mockRepo.AssertExpectations(t)
}

func TestService_FindTasks_RepositoryError(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	query := domain.TaskQuery{UserId: "1", StartDate: &start, EndDate: &end, PageSize: 20}

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", query).Return(domain.TaskPage{}, assert.AnError)

	svc := NewService(mockRepo, testClock)
	_, err := svc.FindTasks(query)

	require.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_GetTaskByID(t *testing.T) {
	now := time.Now().UTC()
	mockRepo := new(mockTaskRepository)
	mockRepo.On("GetByIDAndUserID", 42, "1").Return(domain.Task{
		Id: 42, UserId: "1", Description: "Task 42", CreatedAt: now,
	}, nil)

	svc := NewService(mockRepo, testClock)
	task, err := svc.GetTaskByID("1", 42)

	require.NoError(t, err)
	assert.Equal(t, 42, task.Id)
	assert.Equal(t, "1", task.UserId)
	assert.Equal(t, "Task 42", task.Description)
	mockRepo.AssertExpectations(t)
}

func TestService_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("GetByIDAndUserID", 42, "1").Return(domain.Task{}, domain.ErrTaskNotFound)

	svc := NewService(mockRepo, testClock)
	_, err := svc.GetTaskByID("1", 42)

	assert.ErrorIs(t, err, domain.ErrTaskNotFound)
	mockRepo.AssertExpectations(t)
}
