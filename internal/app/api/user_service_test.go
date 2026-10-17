package api

import (
	"errors"
	"testing"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUserRepository struct {
	findByUsername func(username string) (domain.User, error)
	findByEmail    func(email string) (domain.User, error)
	add            func(user domain.User) (domain.User, error)
}

func (f *fakeUserRepository) FindUserByUsername(username string) (domain.User, error) {
	if f.findByUsername != nil {
		return f.findByUsername(username)
	}
	return domain.User{}, persistence.ErrUserNotFound
}

func (f *fakeUserRepository) FindUserByEmail(email string) (domain.User, error) {
	if f.findByEmail != nil {
		return f.findByEmail(email)
	}
	return domain.User{}, persistence.ErrUserNotFound
}

func (f *fakeUserRepository) Add(user domain.User) (domain.User, error) {
	if f.add != nil {
		return f.add(user)
	}
	return user, nil
}

func validUserRequest() UserRequest {
	return UserRequest{
		Username: "john",
		Email:    "john@example.com",
		Password: "_1ValidPassword",
	}
}

func TestUserServiceCreateUser(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) (domain.User, error) {
			user.Id = 1
			return user, nil
		},
	}

	user, err := NewUserService(repo).createUser(validUserRequest())

	require.NoError(t, err)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "john", user.Username)
}

func TestUserServiceCreateUserUsernameAlreadyExists(t *testing.T) {
	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return domain.User{Id: 1, Username: username}, nil
		},
	}

	_, err := NewUserService(repo).createUser(validUserRequest())

	assert.ErrorIs(t, err, ErrUserAlreadyExists)
}

func TestUserServiceCreateUserEmailAlreadyExists(t *testing.T) {
	repo := &fakeUserRepository{
		findByEmail: func(email string) (domain.User, error) {
			return domain.User{Id: 1, Email: email}, nil
		},
	}

	_, err := NewUserService(repo).createUser(validUserRequest())

	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
}

func TestUserServiceCreateUserRepositoryErrorIsNotConflict(t *testing.T) {
	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return domain.User{}, errors.New("database unavailable")
		},
	}

	_, err := NewUserService(repo).createUser(validUserRequest())

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrUserAlreadyExists)
	assert.NotErrorIs(t, err, ErrEmailAlreadyExists)
}
