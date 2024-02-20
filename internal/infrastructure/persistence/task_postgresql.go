package persistence

import (
	"fmt"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

func NewTaskRepository(db Datasource) domain.TaskRepository {
	return &PostgreRepository{DB: db}
}

func (pg *PostgreRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {

	db, err := pg.connect()
	if err != nil {
		return []domain.Task{}, err
	}

	defer db.Close()

	rows, err := db.Query(`SELECT id, user_id, description, created_at FROM task 
		where created_at >= $1 AND created_at <= $2`, startDate, endDate)

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
