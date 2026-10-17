package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/golang-jwt/jwt/v5"
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

func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) error {
	creds := Credentials{}

	if err := bindJSON(r, &creds); err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if creds.Username == "" || creds.Password == "" {
		return NewHTTPError(http.StatusBadRequest, "username or password is empty")
	}

	user, authorized, err := h.UserService.login(creds.Username, creds.Password)
	if err != nil {
		if errors.Is(err, persistence.ErrUserNotFound) {
			return NewHTTPError(http.StatusUnauthorized, "user not found")
		}
		return NewHTTPError(http.StatusInternalServerError, "failed to login")
	}
	if !authorized {
		return ErrUnauthorized
	}

	now := time.Now().UTC()
	expirationTime := now.Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		UserID: user.Id,
		Name:   user.Username,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(user.Id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	loginResponse := LoginResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
		UserID:    user.Id,
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
