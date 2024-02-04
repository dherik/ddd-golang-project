package app

import (
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

func (s *Service) GetTasks(userId string) []TaskResponse {
	var tasks []domain.Task
	tasks, _ = s.DB.GetTasksFromUser(userId)
	var taskResponses []TaskResponse

	for _, task := range tasks {
		var taskResponse TaskResponse
		copy(task, &taskResponse)
		taskResponses = append(taskResponses, taskResponse)
	}

	return taskResponses
}

func copy(task domain.Task, taskResponse *TaskResponse) TaskResponse {
	return TaskResponse{
		UserId:      task.UserId,
		Description: task.Description,
	}
}
