package api

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type TaskService struct {
	taskRepository domain.TaskRepository
}

func NewTaskService(taskRepository domain.TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) AddTaskToUser(taskRequest TaskRequest) error {
	t, err := toRequest(&taskRequest)
	if err != nil {
		return fmt.Errorf("failed to add task to user: %w", err)
	}

	_, err = s.taskRepository.AddTaskToUser(t.UserId, t)
	if err != nil {
		return fmt.Errorf("failed to add task to user: %w", err)
	}
	return nil
}

func (s *TaskService) FindTasks(startDate time.Time, endDate time.Time) ([]TaskResponse, error) {
	slog.Info(fmt.Sprintf("Find tasks between %s and %s", startDate, endDate))
	tasks, err := s.taskRepository.FindTasks(startDate, endDate)
	if err != nil {
		return []TaskResponse{}, fmt.Errorf("failed finding tasks: %w", err)
	}
	taskResponses := toResponseArray(tasks)
	return taskResponses, nil
}

func (s *TaskService) GetTasksByID(id int) TaskResponse {
	task, _ := s.taskRepository.GetByID(id) //FIXME error
	slog.Info(fmt.Sprintf("Found task with id %d", id))
	return toResponse(task)
}

func toResponseArray(tasks []domain.Task) []TaskResponse {
	taskResponses := []TaskResponse{}
	for _, task := range tasks {
		tr := toResponse(task)
		taskResponses = append(taskResponses, tr)
	}
	return taskResponses
}

func toResponse(task domain.Task) TaskResponse {
	tr := TaskResponse{
		Id:          task.Id,
		UserId:      task.UserId,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
	}
	return tr
}

func toRequest(taskRequest *TaskRequest) (domain.Task, error) {
	task, err := domain.NewTask(taskRequest.UserId, taskRequest.Description)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to convert to request: %w", err)
	}
	return task, nil
}

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) login(username, password string) (bool, error) {
	user, err := s.userRepository.FindUserByUsername(username)
	if err != nil {
		return false, err
	}
	authorized := user.CheckPasswordHash(password)
	return authorized, nil
}
