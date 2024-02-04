package domain

import "github.com/labstack/gommon/log"

// UserRepository defines the contract for managing user data.
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(userID int) error
	FindByID(userID int) (*User, error)
	FindByEmail(email string) (*User, error)
	// Add more methods as needed for user-related data operations
}

// TaskRepository defines the contract for managing task data.
type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(taskID int) error
	FindByID(taskID int) (*Task, error)
	FindAllByUserID(userID int) ([]*Task, error)
	// Add more methods as needed for task-related data operations
}

var m = make(map[int][]*Task, 0)

type MemoryTaskRepository struct {
}

func (sql MemoryTaskRepository) Create(task *Task) error {
	_, ok := m[task.UserId]
	if ok {
		m[task.UserId] = append(m[task.UserId], task)
		log.Info("User cache already exists and element added")
	} else {
		log.Info("User cache create and element added")
		m[task.UserId] = []*Task{task}
	}
	return nil
}

func (sql MemoryTaskRepository) FindAllByUserID(userID int) ([]*Task, error) {
	return m[userID], nil
}
