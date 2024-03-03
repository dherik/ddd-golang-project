package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	TaskService TaskService
}

func NewTaskHandler(taskService TaskService) TaskHandler {
	return TaskHandler{
		TaskService: taskService,
	}
}

func (h *TaskHandler) getTasks(c echo.Context) error {
	startDateParam := c.QueryParam("startDate")
	endDateParam := c.QueryParam("endDate")
	startDate, err := time.Parse(time.RFC3339, startDateParam)
	if err != nil {
		return fmt.Errorf("failed parsing start date: %w", err)
	}
	endDate, err := time.Parse(time.RFC3339, endDateParam)
	if err != nil {
		return fmt.Errorf("failed parsing end date: %w", err)
	}
	tasks, err := h.TaskService.FindTasks(startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed finding tasks: %w", err)
	}
	return c.JSONPretty(http.StatusOK, tasks, "")
}

func (h *TaskHandler) getTaskByID(c echo.Context) error {
	slog.Info("Get all tasks from user")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return fmt.Errorf("failed parsing int using Atoi: %w", err)
	}
	task, err := h.TaskService.GetTasksByID(id)
	if err != nil {
		return fmt.Errorf("failed getting task by id: %w", err)
	}
	return c.JSONPretty(http.StatusOK, task, " ")
}

func (h *TaskHandler) createTask(c echo.Context) error {
	slog.Info("Add new task for user")
	t := TaskRequest{}
	if err := c.Bind(&t); err != nil {
		slog.Error(fmt.Sprintf("failed binding body to task request struct: %s", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the task for the user")
	}

	err := h.TaskService.AddTaskToUser(t)
	if err != nil {
		slog.Error(fmt.Sprintf("failed adding task to user: %s", err.Error()))
		//FIXME return 400 just when is a expected error, otherwise return 500
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the task for the user")
	}

	return c.NoContent(http.StatusCreated)
}
