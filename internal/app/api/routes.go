package api

import (
	"fmt"
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
	UserId      string `json:"userId"`
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
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Routes struct {
	TaskHandler  TaskHandler
	LoginHandler LoginHandler
	UserHandler  UserHandler
}

func NewRouter(taskHandler TaskHandler, loginHandler LoginHandler, userHandler UserHandler) Routes {
	return Routes{
		TaskHandler:  taskHandler,
		LoginHandler: loginHandler,
		UserHandler:  userHandler,
	}
}

func (r *Routes) SetupRoutes(e *echo.Echo) {

	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"), //FIXME read from config
	}

	userGroup := e.Group("/users")
	userGroup.Use(echojwt.WithConfig(config))
	userGroup.POST("", r.UserHandler.createUser)
	// userGroup.GET("", r.UserHandler.getUsers)

	taskGroup := e.Group("/tasks")
	taskGroup.Use(echojwt.WithConfig(config))
	taskGroup.GET("", r.TaskHandler.getTasks)
	taskGroup.GET("/:id", r.TaskHandler.getTaskByID) //FIXME using by user id
	taskGroup.POST("", r.TaskHandler.createTask)

	e.POST("/login", r.LoginHandler.login) //FIXME returning error 500 if user does not exist

	e.GET("/api/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
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
		slog.Error(fmt.Sprint("failed sending custom error response: %w", err))
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
