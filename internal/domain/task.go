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
