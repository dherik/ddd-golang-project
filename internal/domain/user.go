package domain

import (
	"errors"
	"fmt"
	"time"
	"unicode"
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
	Add(user User) (User, error)
}

func NewUser(username, email, plainPassword string) (User, error) {

	user := User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	err := user.Validate()
	if err != nil {
		return User{}, fmt.Errorf("failed creating the user: %w", err)
	}

	password := Password{PlainPassword: plainPassword}
	err = password.Validate()

	if err != nil {
		return User{}, fmt.Errorf("failed creating the password: %w", err)
	}

	password.hashPassword()
	user.Password = password

	return user, nil
}

func hasSpecialCharacter(password string) bool {
	specialChar := false
	for _, char := range password {
		switch char {
		case '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '+', '-', '=':
			specialChar = true
		}
	}
	return specialChar
}

func containsNumber(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func containsUppercaseAndLowercase(s string) bool {
	var hasUppercase, hasLowercase bool
	for _, char := range s {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsLower(char) {
			hasLowercase = true
		}
		if hasUppercase && hasLowercase {
			return true
		}
	}
	return false
}
