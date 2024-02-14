package app

import (
	"fmt"
	"log/slog"
	"time"

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

func (s *Service) FindTasks(startDate time.Time, endDate time.Time) ([]TaskResponse, error) {
	slog.Info(fmt.Sprintf("Find tasks between %s and %s", startDate, endDate))
	tasks, _ := s.DB.FindTasks(startDate, endDate)
	taskResponses := toResponse(tasks)
	return taskResponses, nil
}

func (s *Service) GetTasks(userId string) []TaskResponse {
	var tasks []domain.Task
	tasks, _ = s.DB.GetTasksFromUser(userId)
	slog.Info(fmt.Sprintf("Found %d tasks for user with id %s", len(tasks), userId))
	taskResponses := toResponse(tasks)
	return taskResponses
}

func toResponse(tasks []domain.Task) []TaskResponse {
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
