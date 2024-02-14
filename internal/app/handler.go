package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	router.GET("/tasks", func(c echo.Context) error {
		startDateParam := c.QueryParam("startDate")
		endDateParam := c.QueryParam("endDate")
		startDate, _ := time.Parse(time.RFC3339, startDateParam) //FIXME
		endDate, _ := time.Parse(time.RFC3339, endDateParam)     //FIXME
		tasks, _ := service.FindTasks(startDate, endDate)
		return c.JSONPretty(http.StatusOK, tasks, "")
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

	router.POST("signin", func(c echo.Context) error {
		creds := Credentials{}
		if err := c.Bind(&creds); err != nil {
			slog.Error("Error reading body", slog.String("error", err.Error()))
			return c.JSONPretty(http.StatusBadRequest, nil, "") //FIXME
		}

		expectedPassword, ok := users[creds.Username]

		if !ok || expectedPassword != creds.Password {
			slog.Error("Password invalid")
			return c.JSONPretty(http.StatusUnauthorized, nil, "") //FIXME
		}

		expirationTime := time.Now().Add(1 * time.Hour)
		claims := &Claims{
			Username: creds.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, nil, "") //FIXME
		}

		signInResponse := SignInResponse{
			Token:     tokenString,
			ExpiresAt: expirationTime,
			Username:  creds.Username,
		}
		return c.JSONPretty(http.StatusOK, signInResponse, "")
	})
}

var jwtKey = []byte("my_secret_key")

// FIXME: save in the database
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type SignInResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	// UserID       int       `json:"userId"`
	Username     string   `json:"username"`
	Role         string   `json:"role"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	Metadata     Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
