package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	UserService UserService
	JWTSecret   string
}

func NewLoginHandler(userService UserService, jwtSecret string) LoginHandler {
	return LoginHandler{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

func (h *LoginHandler) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	authorized, err := h.UserService.login(username, password)
	if err != nil {
		return fmt.Errorf("failed login for user %s: %w", username, err)
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
