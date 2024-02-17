package app

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

func (s *TaskService) AddTaskToUser(taskRequest TaskRequest) {
	t := copyRequest(&taskRequest)
	s.taskRepository.AddTaskToUser(t.UserId, t)
}

func (s *TaskService) FindTasks(startDate time.Time, endDate time.Time) ([]TaskResponse, error) {
	slog.Info(fmt.Sprintf("Find tasks between %s and %s", startDate, endDate))
	tasks, _ := s.taskRepository.FindTasks(startDate, endDate)
	taskResponses := toResponse(tasks)
	return taskResponses, nil
}

func (s *TaskService) GetTasks(userId string) []TaskResponse {
	var tasks []domain.Task
	tasks, _ = s.taskRepository.Get(userId)
	slog.Info(fmt.Sprintf("Found %d tasks for user with id %s", len(tasks), userId))
	taskResponses := toResponse(tasks)
	return taskResponses
}

func toResponse(tasks []domain.Task) []TaskResponse {
	taskResponses := []TaskResponse{}
	for _, task := range tasks {
		tr := copyResponse(task)
		taskResponses = append(taskResponses, tr)
	}
	return taskResponses
}

func copyResponse(task domain.Task) TaskResponse {
	tr := TaskResponse{
		Id:          task.Id,
		UserId:      task.UserId,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
	}
	return tr
}

func copyRequest(taskRequest *TaskRequest) domain.Task {
	return domain.NewTask(taskRequest.UserId, taskRequest.Description)
}
