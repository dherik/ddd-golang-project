package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/dherik/ddd-golang-project/internal/app/user"
)

type UserHandler struct {
	UserService *user.Service
}

func NewUserHandler(userService *user.Service) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) error {
	u := UserRequest{}
	if err := bindJSON(r, &u); err != nil {
		slog.Error("failed binding body to user request struct", "err", err)
		return NewHTTPError(http.StatusBadRequest, "failed to create the user")
	}

	err := h.UserService.Create(u.Username, u.Email, u.Password)
	if err != nil {
		if errors.Is(err, user.ErrUsernameTaken) || errors.Is(err, user.ErrEmailTaken) {
			return NewHTTPError(http.StatusConflict, err.Error())
		}

		slog.Error("failed to create user", "err", err)
		return NewHTTPError(http.StatusInternalServerError, "failed to create the user")
	}

	writeNoContent(w, http.StatusCreated)
	return nil
}
