package app

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaskResponse struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	UserId      string `json:"userId"` //FIXME user.id
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
		// tasks := []TaskResponse{
		// 	{"1", "Task 1"},
		// 	{"2", "Task 2"},
		// }
		slog.Info("Get all tasks from user")
		tasks := service.GetTasks(c.Param("id")) //FIXME
		return c.JSONPretty(http.StatusOK, tasks, " ")
	})

	router.POST("/tasks", func(c echo.Context) error {
		// service.AddTask
		slog.Info("Add new task for user")
		t := TaskRequest{}
		if err := c.Bind(&t); err != nil {
			slog.Error("Error reading task body", slog.String("error", err.Error()))
			return err //FIXME
		}

		service.AddTaskToUser(t)

		return c.JSONPretty(http.StatusCreated, nil, "")
	})

	// router.Route("/tasks", func(r chi.Router) {
	// 	r.Get("/", searchTasks)
	// 	r.Route("/{taskId}", func(r chi.Router) {

	// 		r.Get("/", GetTask)

	// 	})
	// })
}

// func getTasks() []TaskResponse {
// 	return nil
// }

func searchTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tasks"))
}

// func GetTask(w http.ResponseWriter, r *http.Request) {
// 	// article := r.Context().Value("article").(*Article)

// 	task := domain.Task{}

// 	if err := render.Render(w, r, NewTaskResponse(&task)); err != nil {
// 		render.Render(w, r, ErrRender(err))
// 		return
// 	}
// }

// func (rd *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
// 	// Pre-processing before a response is marshalled and sent across the wire
// 	rd.Elapsed = 10
// 	return nil
// }

// func NewTaskResponse(task *domain.Task) *TaskResponse {
// 	resp := &TaskResponse{Task: task}

// 	// if resp.User == nil {
// 	// 	if user, _ := dbGetUser(resp.UserID); user != nil {
// 	// 		resp.User = NewUserPayloadResponse(user)
// 	// 	}
// 	// }

// 	return resp
// }

// type TaskResponse struct {
// 	*domain.Task

// 	// User *UserPayload `json:"user,omitempty"`

// 	// We add an additional field to the response here.. such as this
// 	// elapsed computed property
// 	Elapsed int64 `json:"elapsed"`
// }

// func ErrRender(err error) render.Renderer {
// 	return &ErrResponse{
// 		Err:            err,
// 		HTTPStatusCode: 422,
// 		StatusText:     "Error rendering response.",
// 		ErrorText:      err.Error(),
// 	}
// }

// func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
// 	render.Status(r, e.HTTPStatusCode)
// 	return nil
// }

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
