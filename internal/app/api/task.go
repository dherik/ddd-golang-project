package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/task"
	"github.com/dherik/ddd-golang-project/internal/domain"
)

type TaskHandler struct {
	TaskService *task.Service
}

func NewTaskHandler(taskService *task.Service) TaskHandler {
	return TaskHandler{TaskService: taskService}
}

const (
	defaultTaskPageSize = 20
	maxTaskPageSize     = 100
	maxPageOffset       = int(^uint(0) >> 1)
)

func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) error {
	userId, err := authenticatedUserID(r)
	if err != nil {
		return err
	}

	startDate, ok := parseOptionalDate(r.URL.Query().Get("startDate"))
	if !ok {
		return NewHTTPError(http.StatusBadRequest, "invalid start date")
	}
	endDate, ok := parseOptionalDate(r.URL.Query().Get("endDate"))
	if !ok {
		return NewHTTPError(http.StatusBadRequest, "invalid end date")
	}
	limit, ok := parseTaskPageSize(r.URL.Query().Get("limit"))
	if !ok {
		return NewHTTPError(http.StatusBadRequest, "invalid limit")
	}
	offset, ok := parseTaskPageOffset(r.URL.Query().Get("offset"))
	if !ok {
		return NewHTTPError(http.StatusBadRequest, "invalid offset")
	}

	query := domain.TaskQuery{
		UserId:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		PageSize:  limit,
		Offset:    offset,
	}
	page, err := h.TaskService.FindTasks(query)
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, "failed finding tasks")
	}
	return writeJSON(w, http.StatusOK, TasksResponse{
		Items:   toResponseArray(page.Tasks),
		HasMore: page.HasMore,
	})
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request) error {
	userId, err := authenticatedUserID(r)
	if err != nil {
		return err
	}

	slog.Info("get task by id")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return NewHTTPError(http.StatusBadRequest, "invalid task id")
	}
	t, err := h.TaskService.GetTaskByID(userId, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return NewHTTPError(http.StatusNotFound, "task not found")
		}
		return NewHTTPError(http.StatusInternalServerError, "failed getting task by id")
	}
	return writeJSON(w, http.StatusOK, toResponse(t))
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) error {
	userId, err := authenticatedUserID(r)
	if err != nil {
		return err
	}

	slog.Info("add new task for user", "userId", userId)
	t := TaskRequest{}
	if err := bindJSON(r, &t); err != nil {
		slog.Error("failed binding body to task request struct", "err", err)
		return NewHTTPError(http.StatusBadRequest, "failed to create the task for the user")
	}

	err = h.TaskService.CreateTask(userId, t.Description)
	if err != nil {
		if domain.IsTaskValidationError(err) {
			slog.Error("invalid task data", "err", err)
			return NewHTTPError(http.StatusBadRequest, err.Error())
		}
		slog.Error("failed adding task to user", "err", err)
		return NewHTTPError(http.StatusInternalServerError, "failed to create the task for the user")
	}

	writeNoContent(w, http.StatusCreated)
	return nil
}

func toResponseArray(tasks []domain.Task) []TaskResponse {
	taskResponses := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		taskResponses = append(taskResponses, toResponse(t))
	}
	return taskResponses
}

func toResponse(t domain.Task) TaskResponse {
	return TaskResponse{
		Id:          t.Id,
		UserId:      t.UserId,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
	}
}

func parseOptionalDate(raw string) (*time.Time, bool) {
	if raw == "" {
		return nil, true
	}
	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, false
	}
	return &parsed, true
}

func parseTaskPageSize(raw string) (int, bool) {
	if raw == "" {
		return defaultTaskPageSize, true
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > maxTaskPageSize {
		return 0, false
	}
	return limit, true
}

func parseTaskPageOffset(raw string) (int, bool) {
	if raw == "" {
		return 0, true
	}
	offset, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return maxPageOffset, true
		}
		return 0, false
	}
	return int(offset), true
}

func bindJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("empty body")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
