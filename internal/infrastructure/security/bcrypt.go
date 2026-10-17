package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher is the production password-hasher adapter.
type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed hashing password: %w", err)
	}
	return string(hash), nil
}

func (h *BcryptHasher) Verify(hash, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword)) == nil
}
