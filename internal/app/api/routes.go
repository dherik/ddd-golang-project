package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
// See https://github.com/golang-jwt/jwt for more examples
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

func (r *Routes) SetupRoutes(e *echo.Echo) {

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.Recover())
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(r.JWTSecret),
		ContextKey: userContextKey,
	}

	userGroup := e.Group("/users")
	userGroup.Use(echojwt.WithConfig(config))
	userGroup.POST("", r.UserHandler.createUser)
	// userGroup.GET("", r.UserHandler.getUsers)

	taskGroup := e.Group("/tasks")
	taskGroup.Use(echojwt.WithConfig(config))
	taskGroup.GET("", r.TaskHandler.getTasks)
	taskGroup.GET("/:id", r.TaskHandler.getTaskByID)
	taskGroup.POST("", r.TaskHandler.createTask)

	e.POST("/login", r.LoginHandler.login)

	e.GET("/api/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

}

func customHTTPErrorHandler(err error, c echo.Context) {
	// Echo invokes this handler again for errors that RequestLogger already
	// handled (it is configured with HandleError), so a committed response
	// must not be written twice.
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	// Create a custom error response
	errorResponse := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	// Send the custom error response
	err = c.JSON(code, errorResponse)
	if err != nil {
		slog.Error("failed sending custom error response", "err", err)
	}
}

// var jwtKey = []byte("my_secret_key")

// type SignInResponse struct {
// 	Token     string    `json:"token"`
// 	ExpiresAt time.Time `json:"expiresAt"`
// 	// UserID       int       `json:"userId"`
// 	Username     string   `json:"username"`
// 	Role         string   `json:"role"`
// 	RefreshToken string   `json:"refresh_token,omitempty"`
// 	Metadata     Metadata `json:"metadata,omitempty"`
// }

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
