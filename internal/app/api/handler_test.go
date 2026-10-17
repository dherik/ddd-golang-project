package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/task"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTaskRepository struct {
	mock.Mock
}

func (m *mockTaskRepository) GetByIDAndUserID(id int, userId string) (domain.Task, error) {
	args := m.Called(id, userId)
	return args.Get(0).(domain.Task), args.Error(1)
}

func (m *mockTaskRepository) AddTaskToUser(userId string, task domain.Task) error {
	return m.Called(userId, task).Error(0)
}

func (m *mockTaskRepository) FindTasks(query domain.TaskQuery) (domain.TaskPage, error) {
	args := m.Called(query)
	return args.Get(0).(domain.TaskPage), args.Error(1)
}

func withAuthenticatedUser(r *http.Request, userID int) *http.Request {
	claims := &jwtCustomClaims{UserID: userID}
	ctx := context.WithValue(r.Context(), userContextKey, claims)
	return r.WithContext(ctx)
}

func TestTaskHandler_GetTasks(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	now := time.Now().UTC()

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", domain.TaskQuery{
		UserId:    "1",
		StartDate: &start,
		EndDate:   &end,
		PageSize:  20,
	}).Return(domain.TaskPage{
		Tasks: []domain.Task{
			{Id: 1, UserId: "1", Description: "Task 1", CreatedAt: now},
		},
		HasMore: false,
	}, nil)

	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?startDate=2024-01-01T00:00:00Z&endDate=2024-12-31T00:00:00Z", nil)
	rec := httptest.NewRecorder()
	r := withAuthenticatedUser(req, 1)

	err := handler.getTasks(rec, r)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var page TasksResponse
	err = json.Unmarshal(rec.Body.Bytes(), &page)
	require.NoError(t, err)
	assert.Len(t, page.Items, 1)
	assert.False(t, page.HasMore)
	assert.Equal(t, 1, page.Items[0].Id)
	assert.Equal(t, "Task 1", page.Items[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestTaskHandler_GetTasks_InvalidStartDate(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?startDate=invalid&endDate=2024-12-31T00:00:00Z", nil)
	rec := httptest.NewRecorder()
	r := withAuthenticatedUser(req, 1)

	err := handler.getTasks(rec, r)

	require.Error(t, err)
	httpErr, ok := err.(*HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "invalid start date", httpErr.Message)
}

func TestTaskHandler_GetTasks_UnauthorizedWithoutClaims(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?startDate=2024-01-01T00:00:00Z&endDate=2024-12-31T00:00:00Z", nil)
	rec := httptest.NewRecorder()
	r := req

	err := handler.getTasks(rec, r)

	require.Error(t, err)
	httpErr, ok := err.(*HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	now := time.Now().UTC()

	mockRepo := new(mockTaskRepository)
	mockRepo.On("GetByIDAndUserID", 42, "1").Return(domain.Task{
		Id: 42, UserId: "1", Description: "Task 42", CreatedAt: now,
	}, nil)

	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	req.SetPathValue("id", "42")
	rec := httptest.NewRecorder()
	r := withAuthenticatedUser(req, 1)

	err := handler.getTaskByID(rec, r)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var task TaskResponse
	err = json.Unmarshal(rec.Body.Bytes(), &task)
	require.NoError(t, err)
	assert.Equal(t, 42, task.Id)
	assert.Equal(t, "Task 42", task.Description)

	mockRepo.AssertExpectations(t)
}

func TestTaskHandler_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("GetByIDAndUserID", 42, "1").Return(domain.Task{}, domain.ErrTaskNotFound)

	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	req.SetPathValue("id", "42")
	rec := httptest.NewRecorder()
	r := withAuthenticatedUser(req, 1)

	err := handler.getTaskByID(rec, r)

	require.Error(t, err)
	httpErr, ok := err.(*HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestTaskHandler_CreateTaskIgnoresUserIdFromBody(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	mockRepo.On("AddTaskToUser", "1", mock.MatchedBy(func(task domain.Task) bool {
		return task.UserId == "1" && task.Description == "Task from token"
	})).Return(nil)

	svc := task.NewService(mockRepo)
	handler := NewTaskHandler(svc)

	body := strings.NewReader(`{"userId":"2","description":"Task from token"}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r := withAuthenticatedUser(req, 1)

	err := handler.createTask(rec, r)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	mockRepo.AssertExpectations(t)
}
