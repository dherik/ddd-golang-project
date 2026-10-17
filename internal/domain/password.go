package domain

import (
	"errors"
	"strings"
	"unicode"
)

// PasswordHasher hashes and verifies passwords. The domain owns the password
// rules; an adapter owns the algorithm and its cost.
type PasswordHasher interface {
	Hash(plainPassword string) (string, error)
	Verify(hash, plainPassword string) bool
}

type Password struct {
	hash string
}

// NewPassword validates a plaintext password and stores only its hash.
func NewPassword(plainPassword string, hasher PasswordHasher) (Password, error) {
	if err := validatePassword(plainPassword); err != nil {
		return Password{}, err
	}

	hash, err := hasher.Hash(plainPassword)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: hash}, nil
}

// RehydratePassword rebuilds a password value from its stored hash.
func RehydratePassword(hash string) Password {
	return Password{hash: hash}
}

func (p Password) Verify(plainPassword string, hasher PasswordHasher) bool {
	if p.hash == "" {
		return false
	}
	return hasher.Verify(p.hash, plainPassword)
}

func (p Password) Hash() string {
	return p.hash
}

func validatePassword(password string) error {
	if password == "" || len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !hasSpecialCharacter(password) {
		return errors.New("password must contain at least one special character")
	}

	if !containsNumber(password) {
		return errors.New("password must contain at least one number")
	}

	if !containsUppercaseAndLowercase(password) {
		return errors.New("password must contain at least one uppercase letter and one lowercase letter")
	}

	return nil
}

func hasSpecialCharacter(s string) bool {
	return strings.ContainsAny(s, "!@#$%^&*()_+-=")
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
