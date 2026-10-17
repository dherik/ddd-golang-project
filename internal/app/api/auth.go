package api

import (
	"net/http"
	"strconv"
)

const userContextKey = "user"

// authenticatedUserID returns the user ID from the JWT claims stored in the
// request context by the authentication middleware. Handlers must use this
// instead of user IDs supplied by clients.
func authenticatedUserID(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(userContextKey).(*jwtCustomClaims)
	if !ok || claims == nil || claims.UserID <= 0 {
		return "", ErrUnauthorized
	}
	return strconv.Itoa(claims.UserID), nil
}
