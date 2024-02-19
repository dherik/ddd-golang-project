package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func SetupHandler(e *echo.Echo, service *TaskService) {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"), //FIXME
	}

	taskGroup := e.Group("/tasks")
	log.Info(config)
	taskGroup.Use(echojwt.WithConfig(config))

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	taskGroup.GET("", func(c echo.Context) error {
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
		tasks, err := service.FindTasks(startDate, endDate)
		if err != nil {
			return err //FIXME
		}
		return c.JSONPretty(http.StatusOK, tasks, "")
	})

	taskGroup.GET("/:id", func(c echo.Context) error {
		slog.Info("Get all tasks from user")
		tasks := service.GetTasks(c.Param("id")) //FIXME
		return c.JSONPretty(http.StatusOK, tasks, " ")
	})

	taskGroup.POST("", func(c echo.Context) error {
		slog.Info("Add new task for user")
		t := TaskRequest{}
		if err := c.Bind(&t); err != nil {
			slog.Error("Error reading task body", slog.String("error", err.Error()))
			return err //FIXME
		}

		service.AddTaskToUser(t)

		return c.JSONPretty(http.StatusCreated, nil, "")
	})

	// e.POST("/signin", func(c echo.Context) error {
	// 	creds := Credentials{}
	// 	if err := c.Bind(&creds); err != nil {
	// 		slog.Error("Error reading body", slog.String("error", err.Error()))
	// 		return echo.ErrBadRequest
	// 	}

	// 	expectedPassword, ok := users[creds.Username]

	// 	if !ok || expectedPassword != creds.Password {
	// 		slog.Error("Password invalid")
	// 		return echo.ErrUnauthorized
	// 	}

	// 	expirationTime := time.Now().Add(1 * time.Hour)
	// 	claims := &Claims{
	// 		Username: creds.Username,
	// 		RegisteredClaims: jwt.RegisteredClaims{
	// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
	// 		},
	// 	}

	// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 	tokenString, err := token.SignedString(jwtKey)
	// 	if err != nil {
	// 		return echo.ErrInternalServerError
	// 	}

	// 	signInResponse := SignInResponse{
	// 		Token:     tokenString,
	// 		ExpiresAt: expirationTime,
	// 		Username:  creds.Username,
	// 	}
	// 	return c.JSONPretty(http.StatusOK, signInResponse, "")
	// })

	e.POST("/login", login)

	e.GET("/api/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Ok")
	})

}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Jon Snow",
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// var jwtKey = []byte("my_secret_key")

// FIXME: save in the database
// var users = map[string]string{
// 	"user1": "password1",
// 	"user2": "password2",
// }

// type Credentials struct {
// 	Password string `json:"password"`
// 	Username string `json:"username"`
// }

// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.RegisteredClaims
// }

// type SignInResponse struct {
// 	Token     string    `json:"token"`
// 	ExpiresAt time.Time `json:"expiresAt"`
// 	// UserID       int       `json:"userId"`
// 	Username     string   `json:"username"`
// 	Role         string   `json:"role"`
// 	RefreshToken string   `json:"refresh_token,omitempty"`
// 	Metadata     Metadata `json:"metadata,omitempty"`
// }

// type Metadata struct {
// 	Email    string `json:"email"`
// 	FullName string `json:"fullName"`
// }

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
