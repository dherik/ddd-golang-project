package app

import (
	"fmt"
	"log/slog"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

type Service struct {
	DB persistence.DB
}

// func (s Service) GetTasks(userId string) []TaskResponse {
// 	var tasks []TaskResponse
// 	return tasks
// }

type IService interface {
	GetTasks(userId string) []TaskResponse
}

func NewService(db persistence.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) AddTaskToUser(taskRequest TaskRequest) {
	t := copyRequest(&taskRequest)
	s.DB.AddTaskToUser(t.UserId, t)
}

func (s *Service) GetTasks(userId string) []TaskResponse {
	var tasks []domain.Task
	tasks, _ = s.DB.GetTasksFromUser(userId)
	slog.Info(fmt.Sprintf("Found %d tasks for user with id %s", len(tasks), userId))
	var taskResponses []TaskResponse
	for _, task := range tasks {
		tr := copyResponse(task)
		taskResponses = append(taskResponses, tr)
	}

	return taskResponses
}

func copyResponse(task domain.Task) TaskResponse {
	tr := TaskResponse{
		UserId:      task.UserId,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
	}
	return tr
}

func copyRequest(taskRequest *TaskRequest) domain.Task {
	return domain.NewTask(taskRequest.UserId, taskRequest.Description)
}
