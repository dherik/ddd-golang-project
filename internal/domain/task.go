package domain

import "time"

type Task struct {
	Id          int
	UserId      string
	Description string
	CreatedAt   time.Time
}

func NewTask(userId string, description string) Task {
	return Task{
		UserId:      userId,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}
}

type TaskRepository interface {
	GetByID(id int) (Task, error)
	GetByUserID(userId string) ([]Task, error)
	AddTaskToUser(userId string, task Task) (Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]Task, error)
}
