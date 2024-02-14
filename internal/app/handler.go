package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type TaskResponse struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	UserId      string    `json:"userId"` //FIXME user.id
	CreatedAt   time.Time `json:"createdAt"`
}

type TaskRequest struct {
	Description string `json:"description"`
	UserId      string `json:"userId"` //FIXME user.id
}

func SetupHandler(router *echo.Echo, service *Service) {

	router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/tasks/:id", func(c echo.Context) error {
		slog.Info("Get all tasks from user")
		tasks := service.GetTasks(c.Param("id")) //FIXME
		return c.JSONPretty(http.StatusOK, tasks, " ")
	})

	router.POST("/tasks", func(c echo.Context) error {
		slog.Info("Add new task for user")
		t := TaskRequest{}
		if err := c.Bind(&t); err != nil {
			slog.Error("Error reading task body", slog.String("error", err.Error()))
			return err //FIXME
		}

		service.AddTaskToUser(t)

		return c.JSONPretty(http.StatusCreated, nil, "")
	})
}

func searchTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tasks"))
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
