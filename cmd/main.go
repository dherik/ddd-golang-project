package main

import (
	"log"
	"net/http"

	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	// Load application configuration
	cfg, err := config.LoadConfig("config.yaml")

	log.Println(cfg.HTTPPort)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	// db, err := persistence.NewPostgresDB(cfg.DatabaseURL)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize database: %v", err)
	// }
	// defer db.Close()

	// Initialize RabbitMQ connection
	// rabbitMQ, err := infrastructure.NewRabbitMQ(cfg.RabbitMQURL)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	// }
	// defer rabbitMQ.Close()

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root."))
	})

	router.Route("/tasks", func(r chi.Router) {
		r.Get("/", searchTasks)
		r.Route("/{taskId}", func(r chi.Router) {
			// r.Use(TaskCtx)
			r.Get("/", GetTask)
			// r.Delete("/", DeleteTask)
		})
	})

	http.ListenAndServe(":3333", router)

	// Initialize application services
	// appService := app.NewService(db, rabbitMQ)

	// Initialize HTTP server and routes
	// router := app.SetupHTTPRoutes(appService)
	// httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)

	// Start the HTTP server
	// fmt.Printf("Starting HTTP server on port %d...\n", cfg.HTTPPort)
	// if err := http.ListenAndServe(httpAddr, router); err != nil {
	// 	log.Fatalf("HTTP server error: %v", err)
	// }
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

func NewTaskResponse(task *domain.Task) *TaskResponse {
	resp := &TaskResponse{Task: task}

	// if resp.User == nil {
	// 	if user, _ := dbGetUser(resp.UserID); user != nil {
	// 		resp.User = NewUserPayloadResponse(user)
	// 	}
	// }

	return resp
}

func (rd *TaskResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	rd.Elapsed = 10
	return nil
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
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

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
