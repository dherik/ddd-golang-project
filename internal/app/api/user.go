package api

import (
	"errors"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(userService *UserService) UserHandler {
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

	err := h.UserService.createUser(u)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) || errors.Is(err, ErrEmailAlreadyExists) {
			return NewHTTPError(http.StatusConflict, err.Error())
		}

		slog.Error("failed to create user", "err", err)
		return NewHTTPError(http.StatusInternalServerError, "failed to create the user")
	}

	writeNoContent(w, http.StatusCreated)
	return nil
}
