package api

import (
	"log/slog"
	"net/http"
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
		return err //FIXME
	}
	endDate, err := time.Parse(time.RFC3339, endDateParam)
	if err != nil {
		return err //FIXME
	}
	tasks, err := h.TaskService.FindTasks(startDate, endDate)
	if err != nil {
		return err //FIXME
	}
	return c.JSONPretty(http.StatusOK, tasks, "")
}

func (h *TaskHandler) getTaskByID(c echo.Context) error {
	slog.Info("Get all tasks from user")
	tasks := h.TaskService.GetTasks(c.Param("id")) //FIXME
	return c.JSONPretty(http.StatusOK, tasks, " ")
}

func (h *TaskHandler) createTask(c echo.Context) error {
	slog.Info("Add new task for user")
	t := TaskRequest{}
	if err := c.Bind(&t); err != nil {
		slog.Error("Error reading task body", slog.String("error", err.Error()))
		return err //FIXME
	}

	h.TaskService.AddTaskToUser(t)

	return c.NoContent(http.StatusCreated)
}
