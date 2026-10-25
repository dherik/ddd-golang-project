package task

import (
	"errors"
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

type mockTaskPublisher struct {
	mock.Mock
}

func (m *mockTaskPublisher) PublishTaskCreated(event TaskCreated) error {
	args := m.Called(event)
	return args.Error(0)
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

	// The announcement is this test's neighbour, not its subject: every
	// successful save publishes, and publishing is asserted in its own test.
	mockPublisher := new(mockTaskPublisher)
	mockPublisher.On("PublishTaskCreated", mock.Anything).Return(nil)

	err := NewService(mockRepo, testClock, mockPublisher).CreateTask("1", "some description")

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// The fact is built from the task the repository returned: its database id
// (7) exists only on the persisted task, never on the constructed one.
func TestService_CreateTask_PublishesTaskCreatedFact(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.Anything).Return(domain.Task{
		Id:          7,
		UserId:      "1",
		Description: "some description",
		CreatedAt:   testClock.now,
	}, nil)

	mockPublisher := new(mockTaskPublisher)
	mockPublisher.On("PublishTaskCreated", TaskCreated{
		TaskId:      7,
		UserId:      1,
		Description: "some description",
		CreatedAt:   testClock.now,
	}).Return(nil)

	err := NewService(mockRepo, testClock, mockPublisher).CreateTask("1", "some description")

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestService_CreateTask_InvalidDescription(t *testing.T) {
	mockRepo := new(mockTaskRepository)

	err := NewService(mockRepo, testClock, new(mockTaskPublisher)).CreateTask("1", "")

	require.ErrorIs(t, err, domain.ErrDescriptionInvalid)
	mockRepo.AssertNotCalled(t, "AddTaskToUser", mock.Anything, mock.Anything)
}

func TestService_CreateTask_DoesNotPublishWhenValidationFails(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockPublisher := new(mockTaskPublisher)

	err := NewService(mockRepo, testClock, mockPublisher).CreateTask("1", "")

	require.ErrorIs(t, err, domain.ErrDescriptionInvalid)
	mockRepo.AssertNotCalled(t, "AddTaskToUser", mock.Anything, mock.Anything)
	mockPublisher.AssertNotCalled(t, "PublishTaskCreated", mock.Anything)
}

func TestService_CreateTask_SwallowsPublishFailure(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.Anything).Return(domain.Task{
		Id:          7,
		UserId:      "1",
		Description: "some description",
		CreatedAt:   testClock.now,
	}, nil)

	mockPublisher := new(mockTaskPublisher)
	mockPublisher.On("PublishTaskCreated", mock.Anything).Return(errors.New("broker unavailable"))

	err := NewService(mockRepo, testClock, mockPublisher).CreateTask("1", "some description")

	// The save is the truth; the announcement is a side effect (ADR-0005).
	require.NoError(t, err)
	mockPublisher.AssertExpectations(t)
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

	svc := NewService(mockRepo, testClock, new(mockTaskPublisher))
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

	svc := NewService(mockRepo, testClock, new(mockTaskPublisher))
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

	svc := NewService(mockRepo, testClock, new(mockTaskPublisher))
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

	svc := NewService(mockRepo, testClock, new(mockTaskPublisher))
	_, err := svc.GetTaskByID("1", 42)

	assert.ErrorIs(t, err, domain.ErrTaskNotFound)
	mockRepo.AssertExpectations(t)
}
