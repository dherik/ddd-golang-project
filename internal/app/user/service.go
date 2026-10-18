// Package user is the application seam for user operations. It speaks domain
// types only, so no transport drags its own shapes across the seam, and it
// reports failures in its own error vocabulary, so callers map them to their
// own protocol without inspecting domain or adapter errors.
package user

import (
	"errors"
	"fmt"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

var (
	ErrUsernameTaken      = errors.New("a user with the same username already exists")
	ErrEmailTaken         = errors.New("a user with the same email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	userRepository domain.UserRepository
	hasher         domain.PasswordHasher
	clock          domain.Clock
}

func NewService(userRepository domain.UserRepository, hasher domain.PasswordHasher, clock domain.Clock) *Service {
	return &Service{
		userRepository: userRepository,
		hasher:         hasher,
		clock:          clock,
	}
}

// Login authenticates a user. An unknown username and a wrong password both
// return ErrInvalidCredentials so callers cannot tell whether an account exists.
func (s *Service) Login(username, password string) (domain.User, error) {
	found, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, fmt.Errorf("failed to login: %w", err)
	}
	if !found.Password.Verify(password, s.hasher) {
		return domain.User{}, ErrInvalidCredentials
	}
	return found, nil
}

// Create builds a user and hands it to the repository. Conflicts are the
// repository's to report; this module translates them into its own vocabulary
// so callers never inspect adapter errors.
func (s *Service) Create(username, email, plainPassword string) error {
	newUser, err := domain.NewUser(username, email, plainPassword, s.hasher, s.clock)
	if err != nil {
		return fmt.Errorf("failed creating new user: %w", err)
	}

	if err := s.userRepository.Add(newUser); err != nil {
		switch {
		case errors.Is(err, domain.ErrUsernameAlreadyExists):
			return ErrUsernameTaken
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			return ErrEmailTaken
		default:
			return fmt.Errorf("failed saving new user to repository: %w", err)
		}
	}

	return nil
}
