package persistence

import "github.com/dherik/ddd-golang-project/internal/domain"

type Database struct {
	Host string
	Port int
}

func (r *Database) GetTasksFromUser(userId string) ([]domain.Task, error) {
	// implements db interface
	var tasks []domain.Task

	tasks = []domain.Task{
		{UserId: 1, Description: "Task 1"},
		{UserId: 2, Description: "Task 2"},
	}

	return tasks, nil

}

type DB interface {
	GetTasksFromUser(userId string) ([]domain.Task, error)
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTasksFromUser(userId string) ([]domain.Task, error) {
	return r.db.GetTasksFromUser(userId)
}
