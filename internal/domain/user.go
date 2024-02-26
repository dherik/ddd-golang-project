package domain

import (
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
	// Login(username string, password string) (User, error)
	FindUserByUsername(username string) (User, error)
}
