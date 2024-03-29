package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	UserService UserService
	JWTSecret   string
}

// Credentials struct to represent the JSON request body
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginHandler(userService UserService, jwtSecret string) LoginHandler {
	return LoginHandler{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

func (h *LoginHandler) login(c echo.Context) error {

	creds := new(Credentials)

	// Bind the request body to the Credentials struct
	if err := c.Bind(creds); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if creds.Username == "" || creds.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username or password is empty")
	}

	username := creds.Username
	password := creds.Password

	authorized, err := h.UserService.login(username, password)
	if err != nil {
		if errors.Is(err, persistence.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}
		return fmt.Errorf("failed login for user '%s': %w", username, err)
	}
	if !authorized {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	expirationTime := time.Now().Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		username, //FIXME get user's name
		username, //FIXME get user's email
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return fmt.Errorf("failed generating encoded token: %w", err)
	}

	loginResponse := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	return c.JSON(http.StatusOK, loginResponse)
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    int       `json:"userId"`
}
