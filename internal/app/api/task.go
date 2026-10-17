package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	TaskService *TaskService
}

const (
	defaultTaskPageSize = 20
	maxTaskPageSize     = 100
)

func NewTaskHandler(taskService *TaskService) TaskHandler {
	return TaskHandler{
		TaskService: taskService,
	}
}

func (h *TaskHandler) getTasks(c echo.Context) error {
	userId, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	startDate, ok := parseOptionalDate(c.QueryParam("startDate"))
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start date")
	}
	endDate, ok := parseOptionalDate(c.QueryParam("endDate"))
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end date")
	}
	limit, ok := parseTaskPageSize(c.QueryParam("limit"))
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
	}
	offset, ok := parseTaskPageOffset(c.QueryParam("offset"))
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid offset")
	}
	query := domain.TaskQuery{
		UserId:    userId,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Offset:    offset,
	}
	tasks, err := h.TaskService.FindTasks(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed finding tasks")
	}
	return c.JSONPretty(http.StatusOK, tasks, "")
}

func (h *TaskHandler) getTaskByID(c echo.Context) error {
	userId, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	slog.Info("get task by id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}
	task, err := h.TaskService.GetTaskByID(userId, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "task not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed getting task by id")
	}
	return c.JSONPretty(http.StatusOK, task, " ")
}

func (h *TaskHandler) createTask(c echo.Context) error {
	userId, err := authenticatedUserID(c)
	if err != nil {
		return err
	}

	slog.Info("add new task for user", "userId", userId)
	t := TaskRequest{}
	if err := c.Bind(&t); err != nil {
		slog.Error("failed binding body to task request struct", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the task for the user")
	}

	// The authenticated user always owns the new task.
	task, err := domain.NewTask(userId, t.Description)
	if err != nil {
		slog.Error("invalid task data", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.TaskService.AddTaskToUser(task)
	if err != nil {
		slog.Error("failed adding task to user", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create the task for the user")
	}

	return c.NoContent(http.StatusCreated)
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
	offset, err := strconv.Atoi(raw)
	if err != nil || offset < 0 {
		return 0, false
	}
	return offset, true
}
