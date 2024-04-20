package domain

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func (u *User) SetPassword(password string) error {
	pass, err := u.hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed hashing password: %w", err)
	}
	u.Password = pass
	return nil
}

func (u *User) CheckPasswordHash(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	if err != nil {
		slog.Debug("failed comparing hash and password: %w", err)
	}
	return err == nil
}

func (u *User) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hashedPassword), err
}

type UserRepository interface {
	FindUserByUsername(username string) (User, error)
	Add(user User) (User, error)
}

func NewUser(username, email, password string) (User, error) {

	if len(username) > 256 || username == "" {
		return User{}, errors.New("username cannot be empty or longer than 256 characters")
	}

	if len(email) > 256 || email == "" {
		return User{}, errors.New("email cannot be empty or longer than 256 characters")
	}

	if password == "" || len(password) < 8 {
		return User{}, errors.New("password must be at least 8 characters long")
	}

	user := User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	err := user.SetPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("failed creating the password for the new user: %w", err)
	}

	return user, nil
}
