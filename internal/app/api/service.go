package api

import (
	"errors"
	"fmt"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

var (
	ErrUserAlreadyExists  = errors.New("a user with the same username already exists")
	ErrEmailAlreadyExists = errors.New("a user with the same email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) login(username, password string) (domain.User, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, fmt.Errorf("failed to login: %w", err)
	}
	if !user.Password.CheckPasswordHash(password) {
		return domain.User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (s *UserService) createUser(userRequest UserRequest) error {
	newUser, err := domain.NewUser(userRequest.Username, userRequest.Email, userRequest.Password)
	if err != nil {
		return fmt.Errorf("failed creating new user: %w", err)
	}

	if err := s.userRepository.Add(newUser); err != nil {
		switch {
		case errors.Is(err, domain.ErrUsernameAlreadyExists):
			return ErrUserAlreadyExists
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			return ErrEmailAlreadyExists
		default:
			return fmt.Errorf("failed saving new user to repository: %w", err)
		}
	}

	return nil
}
