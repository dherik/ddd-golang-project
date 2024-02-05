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

type RepositoryInterface interface {
	GetTasksFromUser(userId string) ([]domain.Task, error)
}

type Repository struct {
	db DB
}

type MemoryRepository struct {
	memoryDb map[string][]domain.Task
}

func NewRepository(db DB) RepositoryInterface {
	return &Repository{db: db}
}

func NewMemoryRepository() RepositoryInterface {
	return &MemoryRepository{}
}

func (r *Repository) GetTasksFromUser(userId string) ([]domain.Task, error) {
	return r.db.GetTasksFromUser(userId)
}

func (m *MemoryRepository) GetTasksFromUser(userId string) ([]domain.Task, error) {
	return m.memoryDb[userId], nil
}
