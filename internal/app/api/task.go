package api

import (
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

func NewTaskHandler(taskService *TaskService) TaskHandler {
	return TaskHandler{
		TaskService: taskService,
	}
}

func (h *TaskHandler) getTasks(c echo.Context) error {
	startDateParam := c.QueryParam("startDate")
	endDateParam := c.QueryParam("endDate")
	startDate, err := time.Parse(time.RFC3339, startDateParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid start date")
	}
	endDate, err := time.Parse(time.RFC3339, endDateParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid end date")
	}
	tasks, err := h.TaskService.FindTasks(startDate, endDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed finding tasks")
	}
	return c.JSONPretty(http.StatusOK, tasks, "")
}

func (h *TaskHandler) getTaskByID(c echo.Context) error {
	slog.Info("get task by id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}
	task, err := h.TaskService.GetTasksByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed getting task by id")
	}
	return c.JSONPretty(http.StatusOK, task, " ")
}

func (h *TaskHandler) createTask(c echo.Context) error {
	slog.Info("add new task for user")
	t := TaskRequest{}
	if err := c.Bind(&t); err != nil {
		slog.Error("failed binding body to task request struct", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the task for the user")
	}

	task, err := domain.NewTask(t.UserId, t.Description)
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
