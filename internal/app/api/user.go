package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) createUser(c echo.Context) error {
	slog.Info("")

	u := UserRequest{}
	if err := c.Bind(&u); err != nil {
		slog.Error(fmt.Sprintf("failed binding body to user request struct: %s", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the user")
	}

	_, err := h.UserService.createUser(u)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		//FIXME return 400 just when is a expected error, otherwise return 500
		slog.Error(fmt.Sprintf("failed to create user: %s", err.Error()))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create the user")
	}

	return c.NoContent(http.StatusCreated)
}
