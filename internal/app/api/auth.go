package api

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const userContextKey = "user"

// authenticatedUserID returns the user ID from the JWT token stored in the
// request context by the authentication middleware. Handlers must use this
// instead of user IDs supplied by clients.
func authenticatedUserID(c echo.Context) (string, error) {
	token, ok := c.Get(userContextKey).(*jwt.Token)
	if !ok || token == nil {
		return "", echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || claims == nil || claims.UserID <= 0 {
		return "", echo.ErrUnauthorized
	}

	return strconv.Itoa(claims.UserID), nil
}
