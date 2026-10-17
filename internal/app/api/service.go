package api

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

var (
	ErrUserAlreadyExists  = errors.New("a user with the same username already exists")
	ErrEmailAlreadyExists = errors.New("a user with the same email already exists")
)

type TaskService struct {
	taskRepository domain.TaskRepository
}

func NewTaskService(taskRepository domain.TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) AddTaskToUser(task domain.Task) error {
	_, err := s.taskRepository.AddTaskToUser(task.UserId, task)
	if err != nil {
		return fmt.Errorf("failed to add task to user: %w", err)
	}
	return nil
}

func (s *TaskService) FindTasks(userId string, startDate time.Time, endDate time.Time) ([]TaskResponse, error) {
	slog.Info("find tasks between dates", "userId", userId, "startDate", startDate, "endDate", endDate)
	tasks, err := s.taskRepository.FindTasks(userId, startDate, endDate)
	if err != nil {
		return []TaskResponse{}, fmt.Errorf("failed finding tasks: %w", err)
	}
	return toResponseArray(tasks), nil
}

func (s *TaskService) GetTaskByID(userId string, id int) (TaskResponse, error) {
	task, err := s.taskRepository.GetByIDAndUserID(id, userId)
	if err != nil {
		return TaskResponse{}, fmt.Errorf("failed getting task by id: %w", err)
	}
	slog.Info("found task", "id", id, "userId", userId)
	return toResponse(task), nil
}

func toResponseArray(tasks []domain.Task) []TaskResponse {
	taskResponses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		taskResponses = append(taskResponses, toResponse(task))
	}
	return taskResponses
}

func toResponse(task domain.Task) TaskResponse {
	return TaskResponse{
		Id:          task.Id,
		UserId:      task.UserId,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
	}
}

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
