package app

import (
	"net/http"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func SetupHandler(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	router.Route("/tasks", func(r chi.Router) {
		r.Get("/", searchTasks)
		r.Route("/{taskId}", func(r chi.Router) {

			r.Get("/", GetTask)

		})
	})
}

func searchTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tasks"))
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	// article := r.Context().Value("article").(*Article)

	task := domain.Task{}

	if err := render.Render(w, r, NewTaskResponse(&task)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

func (rd *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

func NewTaskResponse(task *domain.Task) *TaskResponse {
	resp := &TaskResponse{Task: task}

	// if resp.User == nil {
	// 	if user, _ := dbGetUser(resp.UserID); user != nil {
	// 		resp.User = NewUserPayloadResponse(user)
	// 	}
	// }

	return resp
}

type TaskResponse struct {
	*domain.Task

	// User *UserPayload `json:"user,omitempty"`

	// We add an additional field to the response here.. such as this
	// elapsed computed property
	Elapsed int64 `json:"elapsed"`
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
