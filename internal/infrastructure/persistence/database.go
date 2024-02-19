package persistence

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	_ "github.com/lib/pq"
)

type Datasource struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (ds *Datasource) ConnectionString() string {
	conn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		ds.Host, ds.Port, ds.User, ds.Password, ds.Name)
	return conn
}

type TaskRepository interface {
	Get(userId string) ([]domain.Task, error)
	AddTaskToUser(userId string, task domain.Task) (domain.Task, error)
	FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error)
}

type PostgreRepository struct {
	DB Datasource
}

func (pg *PostgreRepository) connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", pg.DB.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed connecting to database: %w", err)
	}

	return db, nil
}

func (pg *PostgreRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {

	db, err := pg.connect()
	if err != nil {
		return []domain.Task{}, err
	}

	defer db.Close()

	rows, _ := db.Query(`SELECT id, user_id, description, created_at FROM task 
		where created_at >= $1 AND created_at <= $2`, startDate, endDate)

	tasks := []domain.Task{}
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (pg *PostgreRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {

	db, err := pg.connect()
	if err != nil {
		return domain.Task{}, err
	}
	defer db.Close()

	var id int
	err = db.QueryRow(`INSERT INTO task(user_id, description, created_at)
		VALUES($1, $2, $3) RETURNING id`, task.UserId, task.Description, task.CreatedAt).Scan(&id)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to insert task to database: %w", err)
	}

	task.Id = id // FIXME find task again using id
	return task, nil

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
		m.tasks[userId] = []domain.Task{task}
	}
	return task, nil
}

func NewRepository(db Datasource) TaskRepository {
	return &PostgreRepository{DB: db}
}

func NewMemoryRepository() TaskRepository {
	return &MemoryRepository{
		tasks: make(map[string][]domain.Task),
	}
}

func (pg *PostgreRepository) Get(userId string) ([]domain.Task, error) {

	db, err := pg.connect()
	if err != nil {
		return []domain.Task{}, fmt.Errorf("failed connecting to database: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, user_id, description, created_at FROM task 
		where user_id = $1`, userId)

	if err != nil {
		return []domain.Task{}, fmt.Errorf("failed reading rows: %w", err)
	}

	tasks := []domain.Task{}
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m *MemoryRepository) Get(userId string) ([]domain.Task, error) {
	return m.tasks[userId], nil
}
