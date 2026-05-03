package persistence

import (
	"fmt"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type Task struct {
	Id          int       `db:"id"`
	UserId      string    `db:"user_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"` //TODO rename to explain that is UTC?
	//TODO save local timezone of the task
}

type TaskSqlRepository struct {
	pgsql *PostgreRepository
}

func NewTaskRepository(pgsql *PostgreRepository) domain.TaskRepository {
	return &TaskSqlRepository{pgsql: pgsql}
}

func toTaskDomainList(tasksDb []Task) []domain.Task {
	tasks := make([]domain.Task, 0, len(tasksDb))
	for _, taskDb := range tasksDb {
		tasks = append(tasks, toTaskDomain(taskDb))
	}
	return tasks
}

func toTaskDomain(taskDb Task) domain.Task {
	return domain.Task{
		Id:          taskDb.Id,
		UserId:      taskDb.UserId,
		Description: taskDb.Description,
		CreatedAt:   taskDb.CreatedAt,
	}
}

func (r *TaskSqlRepository) FindTasks(startDate time.Time, endDate time.Time) ([]domain.Task, error) {
	db := r.pgsql.connect()

	rows, err := db.Query(`SELECT id, user_id, description, created_at FROM task 
		where created_at >= $1 AND created_at <= $2`, startDate, endDate)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return toTaskDomainList(tasks), nil
}

func (r *TaskSqlRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	db := r.pgsql.connect()

	var id int
	err := db.QueryRow(`INSERT INTO task(user_id, description, created_at)
		VALUES($1, $2, $3) RETURNING id`, task.UserId, task.Description, task.CreatedAt).Scan(&id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}

	taskEntity, err := r.GetByID(id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to retrieve created task: %w", err)
	}
	return toTaskDomain(Task(taskEntity)), nil
}

func (r *TaskSqlRepository) GetByUserID(userId string) ([]domain.Task, error) {
	db := r.pgsql.connect()

	rows, err := db.Query(`SELECT id, user_id, description, created_at FROM task 
		where user_id = $1`, userId)
	if err != nil {
		return []domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed to scan row: %w", err)
		}
		tasks = append(tasks, task)
	}

	return toTaskDomainList(tasks), nil
}

func (r *TaskSqlRepository) GetByID(id int) (domain.Task, error) {
	db := r.pgsql.connect()

	var task Task
	err := db.QueryRow(`SELECT id, user_id, description, created_at FROM task 
		where id = $1`, id).Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return toTaskDomain(task), nil
}
