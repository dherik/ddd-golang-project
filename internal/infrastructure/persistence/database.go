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

	var tasks = []domain.Task{
		{UserId: "1", Description: "Task 1"},
		{UserId: "2", Description: "Task 2"},
	}

	return tasks, nil

}

type TaskRepository interface {
	Get(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error)
}

type Repository struct {
	db TaskRepository
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
	tasks map[string][]domain.Task
}

// FindTasks implements RepositoryInterface.
func (m *MemoryRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {
	var result []domain.Task
	tasks := m.tasks["1"]
	for _, task := range tasks {
		if task.CreatedAt.After(startDate) && task.CreatedAt.Before(endDate) {
			result = append(result, task)
		}
	}
	return result, nil
}

// AddTaskToUser implements RepositoryInterface.
func (m *MemoryRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	_, ok := m.tasks[userId]
	if ok {
		m.tasks[userId] = append(m.tasks[userId], task)
		slog.Info("User cache already exists and element added")
	} else {
		slog.Info("User cache create and element added")
		// m.tasks = make(map[string][]domain.Task)
		m.tasks[userId] = []domain.Task{task}
	}
	return task, nil
}

func NewRepository(db TaskRepository) TaskRepository {
	return &Repository{db: db}
}

func NewMemoryRepository() TaskRepository {
	return &MemoryRepository{
		tasks: make(map[string][]domain.Task),
	}
}

func (r *Repository) Get(userId string) ([]domain.Task, error) {
	return r.db.Get(userId)
}

func (m *MemoryRepository) Get(userId string) ([]domain.Task, error) {
	return m.tasks[userId], nil
}
