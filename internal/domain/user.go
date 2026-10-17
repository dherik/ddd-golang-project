package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserNotFound          = errors.New("the user was not found")
	ErrUsernameAlreadyExists = errors.New("a user with the same username already exists")
	ErrEmailAlreadyExists    = errors.New("a user with the same email already exists")
)

type User struct {
	Id        int
	Username  string
	Email     string
	Password  Password
	CreatedAt time.Time
}

func (u *User) Validate() error {

	username := u.Username

	if len(username) > 256 || username == "" {
		return errors.New("username cannot be empty or longer than 256 characters")
	}

	email := u.Email

	if len(email) > 256 || email == "" {
		return errors.New("email cannot be empty or longer than 256 characters")
	}

	return nil
}

type UserRepository interface {
	FindUserByUsername(username string) (User, error)
	Add(user User) error
}

func NewUser(username, email, plainPassword string, hasher PasswordHasher, clock Clock) (User, error) {

	user := User{
		Username:  username,
		Email:     email,
		CreatedAt: clock.Now().UTC(),
	}

	err := user.Validate()
	if err != nil {
		return User{}, fmt.Errorf("failed creating the user: %w", err)
	}

	password, err := NewPassword(plainPassword, hasher)
	if err != nil {
		return User{}, fmt.Errorf("failed creating the password: %w", err)
	}

	user.Password = password

	return user, nil
}
