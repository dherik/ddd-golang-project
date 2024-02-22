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
}

func NewLoginHandler(userService UserService) LoginHandler {
	return LoginHandler{
		UserService: userService,
	}
}

func (h *LoginHandler) login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	// if username != "jon" || password != "shhh!" {
	// 	return echo.ErrUnauthorized
	// }

	authorized, err := h.UserService.login(username, password)
	if err != nil {
		return fmt.Errorf("failed login for user %s: %w", username, err) //FIXME
	}
	if !authorized {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	expirationTime := time.Now().Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		"Jon Snow",        //FIXME
		"email@email.com", //FIXME
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err //FIXME
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
