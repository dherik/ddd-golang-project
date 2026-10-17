package domain

import (
	"errors"
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

func NewTask(userId string, description string, clock Clock) (Task, error) {

	if len(description) == 0 {
		return Task{}, ErrDescriptionInvalid
	}

	if userId == "" {
		return Task{}, ErrUserIdInvalid
	}

	if len(description) > 1024 {
		return Task{}, ErrDescriptionTooLong
	}

	return Task{
		UserId:      userId,
		Description: description,
		CreatedAt:   clock.Now().UTC(),
	}, nil
}

type TaskQuery struct {
	UserId    string
	StartDate *time.Time
	EndDate   *time.Time
	PageSize  int
	Offset    int
}

type TaskPage struct {
	Tasks   []Task
	HasMore bool
}

type TaskRepository interface {
	GetByIDAndUserID(id int, userId string) (Task, error)
	AddTaskToUser(userId string, task Task) error
	FindTasks(query TaskQuery) (TaskPage, error)
}
