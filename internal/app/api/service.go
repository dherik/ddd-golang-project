package api

import (
	"errors"
	"fmt"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

var (
	ErrUserAlreadyExists  = errors.New("a user with the same username already exists")
	ErrEmailAlreadyExists = errors.New("a user with the same email already exists")
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) login(username, password string) (domain.User, bool, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		return domain.User{}, false, fmt.Errorf("failed to login: %w", err)
	}
	return user, user.Password.CheckPasswordHash(password), nil
}

func (s *UserService) createUser(userRequest UserRequest) (domain.User, error) {
	newUser, err := domain.NewUser(userRequest.Username, userRequest.Email, userRequest.Password)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed creating new user: %w", err)
	}

	if err := s.checkUserDoesNotExist(newUser); err != nil {
		return domain.User{}, err
	}

	user, err := s.userRepository.Add(newUser)
	if err != nil {
		if errors.Is(err, persistence.ErrUsernameAlreadyExists) {
			return domain.User{}, ErrUserAlreadyExists
		}
		if errors.Is(err, persistence.ErrEmailAlreadyExists) {
			return domain.User{}, ErrEmailAlreadyExists
		}
		return domain.User{}, fmt.Errorf("failed saving new user to repository: %w", err)
	}
	return user, nil
}

func (s *UserService) checkUserDoesNotExist(user domain.User) error {
	_, err := s.userRepository.FindUserByUsername(user.Username)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, persistence.ErrUserNotFound) {
		return fmt.Errorf("failed checking username availability: %w", err)
	}

	_, err = s.userRepository.FindUserByEmail(user.Email)
	if err == nil {
		return ErrEmailAlreadyExists
	}
	if !errors.Is(err, persistence.ErrUserNotFound) {
		return fmt.Errorf("failed checking email availability: %w", err)
	}

	return nil
}
