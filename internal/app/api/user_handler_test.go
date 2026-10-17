package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandlerCreateUserRepositoryFailureReturns500(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return errors.New("database unavailable")
		},
	}
	handler := NewUserHandler(NewUserService(repo, fakePasswordHasher{}, testClock))

	body := `{"username":"john","email":"john@example.com","password":"_1ValidPassword"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	rec := httptest.NewRecorder()

	err := handler.createUser(rec, req)

	require.Error(t, err)
	httpErr, ok := err.(*HTTPError)
	require.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
}
