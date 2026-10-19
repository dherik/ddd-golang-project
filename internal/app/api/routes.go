package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TaskResponse struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	UserId      string    `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TaskRequest struct {
	Description string `json:"description"`
}

type TasksResponse struct {
	Items   []TaskResponse `json:"items"`
	HasMore bool           `json:"hasMore"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	UserID int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Routes struct {
	TaskHandler  TaskHandler
	LoginHandler LoginHandler
	UserHandler  UserHandler
	JWTSecret    string
}

func NewRouter(taskHandler TaskHandler, loginHandler LoginHandler, userHandler UserHandler, jwtSecret string) Routes {
	return Routes{
		TaskHandler:  taskHandler,
		LoginHandler: loginHandler,
		UserHandler:  userHandler,
		JWTSecret:    jwtSecret,
	}
}

// SetupRoutes registers all handlers on the given mux.
func (r *Routes) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ok"))
	})

	mux.HandleFunc("POST /login", handle(r.LoginHandler.login))

	mux.HandleFunc("POST /users", handle(requireAuth(r.JWTSecret, r.UserHandler.createUser)))
	mux.HandleFunc("GET /tasks", handle(requireAuth(r.JWTSecret, r.TaskHandler.getTasks)))
	mux.HandleFunc("GET /tasks/{id}", handle(requireAuth(r.JWTSecret, r.TaskHandler.getTaskByID)))
	mux.HandleFunc("POST /tasks", handle(requireAuth(r.JWTSecret, r.TaskHandler.createTask)))
}

// handle adapts an error-returning handler to http.HandlerFunc.
func handle(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			CustomHTTPErrorHandler(err, w, r)
		}
	}
}

// requireAuth wraps a handler with JWT verification.
func requireAuth(jwtSecret string, next func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		auth := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			return ErrUnauthorized
		}
		tokenStr := strings.TrimSpace(auth[len(prefix):])
		if tokenStr == "" {
			return ErrUnauthorized
		}

		token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return ErrUnauthorized
		}
		claims, ok := token.Claims.(*jwtCustomClaims)
		if !ok || claims == nil || claims.UserID <= 0 {
			return ErrUnauthorized
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		return next(w, r.WithContext(ctx))
	}
}

// HTTPError carries an HTTP status code and message, mirroring Echo's.
type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string { return e.Message }

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}

var ErrUnauthorized = NewHTTPError(http.StatusUnauthorized, "Unauthorized")

func CustomHTTPErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		message = he.Message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"timestamp": time.Now().Format(time.RFC3339),
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	})
}

// writeJSON writes v as compact JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// writeNoContent writes an empty response with the given status code.
func writeNoContent(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}
