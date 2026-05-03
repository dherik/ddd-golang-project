package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	UserService *UserService
	JWTSecret   string
}

func NewLoginHandler(userService *UserService, jwtSecret string) LoginHandler {
	return LoginHandler{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

func (h *LoginHandler) login(c echo.Context) error {
	creds := new(Credentials)

	if err := c.Bind(creds); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if creds.Username == "" || creds.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username or password is empty")
	}

	authorized, err := h.UserService.login(creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, persistence.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to login")
	}
	if !authorized {
		return echo.ErrUnauthorized
	}

	expirationTime := time.Now().Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		creds.Username, //FIXME get user's name
		creds.Username, //FIXME get user's email
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	loginResponse := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	return c.JSON(http.StatusOK, loginResponse)
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    int       `json:"userId"`
}
