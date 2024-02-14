package persistence

import (
	"fmt"
	"log/slog"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type Database struct {
	Host string
	Port int
}

func (r *Database) GetTasksFromUser(userId string) ([]domain.Task, error) {
	// implements db interface
	var tasks []domain.Task

	tasks = []domain.Task{
		{UserId: "1", Description: "Task 1"},
		{UserId: "2", Description: "Task 2"},
	}

	return tasks, nil

}

type DB interface {
	GetTasksFromUser(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
}

type RepositoryInterface interface {
	GetTasksFromUser(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
}

type Repository struct {
	db DB
}

// AddTaskToUser implements RepositoryInterface.
func (*Repository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	panic("unimplemented")
}

type MemoryRepository struct {
	memoryDb map[string][]domain.Task
}

// AddTaskToUser implements RepositoryInterface.
func (m *MemoryRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	_, ok := m.memoryDb[userId]
	if ok {
		m.memoryDb[userId] = append(m.memoryDb[userId], task)
		slog.Info("User cache already exists and element added")
	} else {
		slog.Info("User cache create and element added")
		m.memoryDb = make(map[string][]domain.Task)
		m.memoryDb[userId] = []domain.Task{task}
	}
	return task, nil
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
	slog.Info(fmt.Sprintf("MemoryDB: %s", m.memoryDb))
	slog.Info(fmt.Sprintf("UserId: %s", userId))
	return m.memoryDb[userId], nil
}
