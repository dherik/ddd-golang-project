package domain

import (
	"errors"
	"log/slog"
	"time"
)

var ErrDescriptionInvalid = errors.New("the description is not valid")
var ErrDescriptionTooLong = errors.New("the description is too long")
var ErrUserIdInvalid = errors.New("the user id is not valid")
var ErrTaskNotFound = errors.New("the task was not found")

type Task struct {
	Id          int
	UserId      string
	Description string
	CreatedAt   time.Time
}

func NewTask(userId string, description string) (Task, error) {

	if len(description) == 0 {
		slog.Error("invalid description: value is empty")
		return Task{}, ErrDescriptionInvalid
	}

	if userId == "" {
		slog.Error("invalid user id", "userId", userId)
		return Task{}, ErrUserIdInvalid
	}

	if len(description) > 1024 {
		slog.Error("invalid description: value is too long", "description", description)
		return Task{}, ErrDescriptionTooLong
	}

	return Task{
		UserId:      userId,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

type TaskQuery struct {
	UserId    string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

type TaskPage struct {
	Tasks   []Task
	HasMore bool
}

type TaskRepository interface {
	GetByIDAndUserID(id int, userId string) (Task, error)
	AddTaskToUser(userId string, task Task) (Task, error)
	FindTasks(query TaskQuery) (TaskPage, error)
}
