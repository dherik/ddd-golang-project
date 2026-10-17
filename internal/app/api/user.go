package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *UserService
}

func NewUserHandler(userService *UserService) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) createUser(c echo.Context) error {
	u := UserRequest{}
	if err := c.Bind(&u); err != nil {
		slog.Error("failed binding body to user request struct", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the user")
	}

	_, err := h.UserService.createUser(u)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) || errors.Is(err, ErrEmailAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		slog.Error("failed to create user", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create the user")
	}

	return c.NoContent(http.StatusCreated)
}
