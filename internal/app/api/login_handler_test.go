package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dherik/ddd-golang-project/internal/app/user"
	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	stored, err := domain.NewUser("john", "john@example.com", "_1ValidPassword", fakePasswordHasher{}, testClock)
	require.NoError(t, err)

	tests := []struct {
		name string
		repo *fakeUserRepository
	}{
		{
			name: "unknown user",
			repo: &fakeUserRepository{},
		},
		{
			name: "wrong password",
			repo: &fakeUserRepository{
				findByUsername: func(username string) (domain.User, error) {
					return stored, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewLoginHandler(user.NewService(tt.repo, fakePasswordHasher{}, testClock), "test-secret")

			body := `{"username":"john","password":"some_password"}`
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
			rec := httptest.NewRecorder()

			err := handler.login(rec, req)

			require.Error(t, err)
			httpErr, ok := err.(*HTTPError)
			require.True(t, ok)
			assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
			assert.Equal(t, "invalid credentials", httpErr.Message)
		})
	}
}

func TestLoginHandlerReturnsTokenForValidCredentials(t *testing.T) {
	stored, err := domain.NewUser("john", "john@example.com", "_1ValidPassword", fakePasswordHasher{}, testClock)
	require.NoError(t, err)

	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return stored, nil
		},
	}
	handler := NewLoginHandler(user.NewService(repo, fakePasswordHasher{}, testClock), "test-secret")

	body := `{"username":"john","password":"_1ValidPassword"}`
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	err = handler.login(rec, req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response LoginResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, stored.Id, response.UserID)
}
