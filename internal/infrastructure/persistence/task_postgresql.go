package persistence

import (
	"fmt"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type TaskSqlRepository struct {
	pgsql PostgreRepository
}

func NewTaskRepository(pgsql PostgreRepository) domain.TaskRepository {
	return &TaskSqlRepository{pgsql: pgsql}
}

func (r *TaskSqlRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {

	db, err := r.pgsql.connect()
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

func (r *TaskSqlRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {

	db, err := r.pgsql.connect()
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

	return r.GetByID(id)

	// task.Id = id // FIXME find task again using id
	// return task, nil

}

func (r *TaskSqlRepository) GetByUserID(userId string) ([]domain.Task, error) {

	db, err := r.pgsql.connect()
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

func (r *TaskSqlRepository) GetByID(id int) (domain.Task, error) {

	db, err := r.pgsql.connect()
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed connecting to database: %w", err)
	}
	defer db.Close()

	var task domain.Task
	err = db.QueryRow(`SELECT id, user_id, description, created_at FROM task 
		where id = $1`, id).Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed reading rows: %w", err)
	}

	return task, nil

	// tasks := []domain.Task{}
	// for rows.Next() {
	// 	var task domain.Task
	// 	err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
	// 	if err != nil {
	// 		return []domain.Task{}, fmt.Errorf("failed scanning row: %w", err)
	// 	}
	// 	tasks = append(tasks, task)
	// }

	// return tasks, nil
}
