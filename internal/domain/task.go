package domain

import "time"

type Task struct {
	Id          int       `db:"id"` //FIXME create persistence.Task
	UserId      string    `db:"user_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewTask(userId string, desc string) Task {
	return Task{
		UserId:      userId,
		Description: desc,
		CreatedAt:   time.Now().UTC(),
	}
}

type TaskRepository interface {
	GetByID(id int) (Task, error)
	GetByUserID(userId string) ([]Task, error)
	AddTaskToUser(userId string, task Task) (Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]Task, error)
}
