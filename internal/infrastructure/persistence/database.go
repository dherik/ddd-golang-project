package persistence

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

type DatabaseConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// func (r *DatabaseConnection) GetTasksFromUser(userId string) ([]domain.Task, error) {

// 	var tasks = []domain.Task{
// 		{UserId: "1", Description: "Task 1"},
// 		{UserId: "2", Description: "Task 2"},
// 	}

// 	return tasks, nil

// }

type TaskRepository interface {
	Get(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error)
}

type PostgreRepository struct {
	DB DatabaseConnection
}

func (pg *PostgreRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "postgres"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, pg.DB.Port, pg.DB.User, pg.DB.Password, pg.DB.Name)

	log.Info(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return []domain.Task{}, nil
}

func (*PostgreRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	panic("unimplemented")
}

type MemoryRepository struct {
	tasks map[string][]domain.Task
}

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

func NewRepository(db DatabaseConnection) TaskRepository {
	return &PostgreRepository{DB: db}
}

func NewMemoryRepository() TaskRepository {
	return &MemoryRepository{
		tasks: make(map[string][]domain.Task),
	}
}

func (r *PostgreRepository) Get(userId string) ([]domain.Task, error) {
	// return r.db.Get(userId)
	return nil, nil
}

func (m *MemoryRepository) Get(userId string) ([]domain.Task, error) {
	return m.tasks[userId], nil
}
