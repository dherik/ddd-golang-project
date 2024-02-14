package persistence

import (
	"log/slog"
	"time"

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
	FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error)
}

type RepositoryInterface interface {
	GetTasksFromUser(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error)
}

type Repository struct {
	db DB
}

// FindTasks implements RepositoryInterface.
func (*Repository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {
	panic("unimplemented")
}

// AddTaskToUser implements RepositoryInterface.
func (*Repository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	panic("unimplemented")
}

type MemoryRepository struct {
	memoryDb map[string][]domain.Task
}

// FindTasks implements RepositoryInterface.
func (m *MemoryRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {
	var result []domain.Task
	tasks := m.memoryDb["1"]
	for _, task := range tasks {
		if task.CreatedAt.After(startDate) && task.CreatedAt.Before(endDate) {
			result = append(result, task)
		}
	}
	return result, nil
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
	return m.memoryDb[userId], nil
}
