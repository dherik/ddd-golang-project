package domain

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

var ErrDescriptionInvalid = errors.New("the description is not valid")
var ErrDescriptionTooLong = errors.New("the description is too long")
var ErrUserIdInvalid = errors.New("the user id is not valid")

type Task struct {
	Id          int
	UserId      string
	Description string
	CreatedAt   time.Time
}

func NewTask(userId string, description string) (Task, error) {

	if len(description) == 0 {
		slog.Error(fmt.Sprintf("invalid description %s: value is empty", description))
		return Task{}, ErrDescriptionInvalid
	}

	if userId == "" {
		slog.Error(fmt.Sprintf("invalid user id %s", userId))
		return Task{}, ErrUserIdInvalid
	}

	if len(description) > 1024 {
		slog.Error(fmt.Sprintf("invalid description %s: value is too long", description))
		return Task{}, ErrDescriptionTooLong
	}

	return Task{
		UserId:      userId,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

type TaskRepository interface {
	GetByID(id int) (Task, error)
	GetByUserID(userId string) ([]Task, error)
	AddTaskToUser(userId string, task Task) (Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]Task, error)
}
