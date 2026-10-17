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

func TestUserHandlerCreateUserConflictReturns409(t *testing.T) {
	tests := []struct {
		name        string
		repoErr     error
		wantMessage string
	}{
		{
			name:        "username conflict",
			repoErr:     domain.ErrUsernameAlreadyExists,
			wantMessage: "a user with the same username already exists",
		},
		{
			name:        "email conflict",
			repoErr:     domain.ErrEmailAlreadyExists,
			wantMessage: "a user with the same email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepository{
				add: func(user domain.User) error {
					return tt.repoErr
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
			assert.Equal(t, http.StatusConflict, httpErr.Code)
			assert.Equal(t, tt.wantMessage, httpErr.Message)
		})
	}
}
