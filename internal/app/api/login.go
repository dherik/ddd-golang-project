package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/user"
	"github.com/golang-jwt/jwt/v5"
)

type LoginHandler struct {
	UserService *user.Service
	JWTSecret   string
}

func NewLoginHandler(userService *user.Service, jwtSecret string) LoginHandler {
	return LoginHandler{
		UserService: userService,
		JWTSecret:   jwtSecret,
	}
}

func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) error {
	creds := Credentials{}

	if err := bindJSON(r, &creds); err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if creds.Username == "" || creds.Password == "" {
		return NewHTTPError(http.StatusBadRequest, "username or password is empty")
	}

	loggedIn, err := h.UserService.Login(creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			return NewHTTPError(http.StatusUnauthorized, user.ErrInvalidCredentials.Error())
		}
		return NewHTTPError(http.StatusInternalServerError, "failed to login")
	}

	now := time.Now().UTC()
	expirationTime := now.Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		UserID:    loggedIn.Id,
		Name:      loggedIn.Username,
		Email:     loggedIn.Email,
		Subject:   strconv.Itoa(loggedIn.Id),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	loginResponse := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
		UserID:    loggedIn.Id,
	}

	return writeJSON(w, http.StatusOK, loginResponse)
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
