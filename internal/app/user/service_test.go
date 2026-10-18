package user

import (
	"errors"
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUsername = "john"
	testEmail    = "john@example.com"
	testPassword = "_1ValidPassword"
)

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(plainPassword string) (string, error) {
	return "hashed:" + plainPassword, nil
}

func (fakePasswordHasher) Verify(hash, plainPassword string) bool {
	return hash == "hashed:"+plainPassword
}

var testClock = fixedClock{now: time.Date(2026, 10, 17, 12, 0, 0, 0, time.UTC)}

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

func TestServiceCreate(t *testing.T) {
	var added domain.User
	findCalled := false
	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			findCalled = true
			return domain.User{}, domain.ErrUserNotFound
		},
		add: func(user domain.User) error {
			added = user
			return nil
		},
	}

	err := NewService(repo, fakePasswordHasher{}, testClock).Create(testUsername, testEmail, testPassword)

	require.NoError(t, err)
	assert.False(t, findCalled, "Create must not pre-check the repository")
	assert.Equal(t, testUsername, added.Username)
	assert.Equal(t, testEmail, added.Email)
	assert.Equal(t, "hashed:"+testPassword, added.Password.Hash())
	assert.True(t, added.CreatedAt.Equal(testClock.now))
}

func TestServiceCreateUsernameConflictOnInsert(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return domain.ErrUsernameAlreadyExists
		},
	}

	err := NewService(repo, fakePasswordHasher{}, testClock).Create(testUsername, testEmail, testPassword)

	assert.ErrorIs(t, err, ErrUsernameTaken)
}

func TestServiceCreateEmailConflictOnInsert(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return domain.ErrEmailAlreadyExists
		},
	}

	err := NewService(repo, fakePasswordHasher{}, testClock).Create(testUsername, testEmail, testPassword)

	assert.ErrorIs(t, err, ErrEmailTaken)
}

func TestServiceCreateRepositoryErrorIsNotConflict(t *testing.T) {
	repo := &fakeUserRepository{
		add: func(user domain.User) error {
			return errors.New("database unavailable")
		},
	}

	err := NewService(repo, fakePasswordHasher{}, testClock).Create(testUsername, testEmail, testPassword)

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrUsernameTaken)
	assert.NotErrorIs(t, err, ErrEmailTaken)
}

func TestServiceLogin(t *testing.T) {
	stored, err := domain.NewUser(testUsername, testEmail, testPassword, fakePasswordHasher{}, testClock)
	require.NoError(t, err)

	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return stored, nil
		},
	}

	got, err := NewService(repo, fakePasswordHasher{}, testClock).Login(testUsername, testPassword)

	require.NoError(t, err)
	assert.Equal(t, testUsername, got.Username)
}

func TestServiceLoginUnknownUserReturnsInvalidCredentials(t *testing.T) {
	repo := &fakeUserRepository{}

	_, err := NewService(repo, fakePasswordHasher{}, testClock).Login("missing", testPassword)

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestServiceLoginWrongPasswordReturnsInvalidCredentials(t *testing.T) {
	stored, err := domain.NewUser(testUsername, testEmail, testPassword, fakePasswordHasher{}, testClock)
	require.NoError(t, err)

	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return stored, nil
		},
	}

	_, err = NewService(repo, fakePasswordHasher{}, testClock).Login(testUsername, "another_password")

	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestServiceLoginRepositoryErrorIsNotInvalidCredentials(t *testing.T) {
	repo := &fakeUserRepository{
		findByUsername: func(username string) (domain.User, error) {
			return domain.User{}, errors.New("database unavailable")
		},
	}

	_, err := NewService(repo, fakePasswordHasher{}, testClock).Login(testUsername, testPassword)

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrInvalidCredentials)
}
