package api

import (
	"errors"
	"testing"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUserRepository struct {
	findByUsername func(username string) (domain.User, error)
	add            func(user domain.User) error
}

func (f *fakeUserRepository) FindUserByUsername(username string) (domain.User, error) {
	if f.findByUsername != nil {
		return f.findByUsername(username)
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (f *fakeUserRepository) Add(user domain.User) error {
	if f.add != nil {
		return f.add(user)
	}
	return nil
}

func validUserRequest() UserRequest {
	return UserRequest{
		Username: "john",
		Email:    "john@example.com",
		Password: "_1ValidPassword",
	}
}

func TestUserServiceCreateUser(t *testing.T) {
	var added domain.User
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			added = user
			return nil
		},
	}

	err := NewUserService(repo).createUser(validUserRequest())

	require.NoError(t, err)
	assert.Equal(t, "john", added.Username)
	assert.Equal(t, "john@example.com", added.Email)
	assert.NotEmpty(t, added.Password.HashPassword)
}

func TestUserServiceCreateUserUsernameConflictOnInsert(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return domain.ErrUsernameAlreadyExists
		},
	}

	err := NewUserService(repo).createUser(validUserRequest())

	assert.ErrorIs(t, err, ErrUserAlreadyExists)
}

func TestUserServiceCreateUserEmailConflictOnInsert(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return domain.ErrEmailAlreadyExists
		},
	}

	err := NewUserService(repo).createUser(validUserRequest())

	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestUserServiceCreateUserRepositoryErrorIsNotConflict(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return errors.New("database unavailable")
		},
	}

	err := NewUserService(repo).createUser(validUserRequest())

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrUserAlreadyExists)
	assert.NotErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestUserServiceLogin(t *testing.T) {
	user, err := domain.NewUser("john", "john@example.com", "_1ValidPassword")
	require.NoError(t, err)

	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return user, nil
		},
	}

	got, err := NewUserService(repo).login("john", "_1ValidPassword")

	require.NoError(t, err)
	assert.Equal(t, "john", got.Username)
}

func TestUserServiceLoginUnknownUserReturnsInvalidCredentials(t *testing.T) {
	repo := &fakeUserRepository{}

	_, err := NewUserService(repo).login("missing", "some_password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestUserServiceLoginWrongPasswordReturnsInvalidCredentials(t *testing.T) {
	user, err := domain.NewUser("john", "john@example.com", "_1ValidPassword")
	require.NoError(t, err)

	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return user, nil
		},
	}

	_, err = NewUserService(repo).login("john", "another_password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestUserServiceLoginRepositoryErrorIsNotInvalidCredentials(t *testing.T) {
	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return domain.User{}, errors.New("database unavailable")
		},
	}

	_, err := NewUserService(repo).login("john", "_1ValidPassword")

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrInvalidCredentials)
}
