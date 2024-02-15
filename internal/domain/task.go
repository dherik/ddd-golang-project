package domain

import "time"

type Task struct {
	UserId      string
	Description string
	CreatedAt   time.Time
}

func NewTask(userId string, desc string) Task {
	return Task{
		UserId:      userId,
		Description: desc,
		CreatedAt:   time.Now().UTC(),
	}
}

type TaskRepository interface {
	Get(userId string) ([]Task, error)
	AddTaskToUser(userId string, task Task) (Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]Task, error)
}
