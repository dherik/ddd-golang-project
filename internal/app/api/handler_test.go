package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler_GetTasks(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	now := time.Now().UTC()

	mockRepo := new(mockTaskRepository)
	mockRepo.On("FindTasks", start, end).Return([]domain.Task{
		{Id: 1, UserId: "1", Description: "Task 1", CreatedAt: now},
	}, nil)

	svc := NewTaskService(mockRepo)
	handler := NewTaskHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks?startDate=2024-01-01T00:00:00Z&endDate=2024-12-31T00:00:00Z", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.getTasks(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var tasks []TaskResponse
	err = json.Unmarshal(rec.Body.Bytes(), &tasks)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, 1, tasks[0].Id)
	assert.Equal(t, "Task 1", tasks[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestTaskHandler_GetTasks_InvalidStartDate(t *testing.T) {
	mockRepo := new(mockTaskRepository)
	svc := NewTaskService(mockRepo)
	handler := NewTaskHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks?startDate=invalid&endDate=2024-12-31T00:00:00Z", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.getTasks(c)

	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.Equal(t, "invalid start date", httpErr.Message)
}

func TestTaskHandler_GetTaskByID(t *testing.T) {
	now := time.Now().UTC()

	mockRepo := new(mockTaskRepository)
	mockRepo.On("GetByID", 42).Return(domain.Task{
		Id: 42, UserId: "1", Description: "Task 42", CreatedAt: now,
	}, nil)

	svc := NewTaskService(mockRepo)
	handler := NewTaskHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("42")

	err := handler.getTaskByID(c)
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
	mockRepo.On("GetByID", 42).Return(domain.Task{}, domain.ErrTaskNotFound)

	svc := NewTaskService(mockRepo)
	handler := NewTaskHandler(svc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks/42", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("42")

	err := handler.getTaskByID(c)

	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)

	mockRepo.AssertExpectations(t)
}
