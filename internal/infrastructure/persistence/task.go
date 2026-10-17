package persistence

import (
	"database/sql"
	"errors"
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

func (r *TaskSqlRepository) FindTasks(query domain.TaskQuery) (domain.TaskPage, error) {
	if query.Limit < 1 {
		return domain.TaskPage{}, fmt.Errorf("page size must be at least 1")
	}
	if query.Offset < 0 {
		return domain.TaskPage{}, fmt.Errorf("offset must not be negative")
	}

	db := r.pgsql.connect()

	// Fetch one row beyond the page so HasMore is known without a COUNT.
	rows, err := db.Query(`SELECT id, user_id, description, created_at FROM task
		where user_id = $1
		  AND ($2::timestamptz IS NULL OR created_at >= $2)
		  AND ($3::timestamptz IS NULL OR created_at <= $3)
		ORDER BY created_at DESC, id DESC
		LIMIT $4 OFFSET $5`, query.UserId, query.StartDate, query.EndDate, query.Limit+1, query.Offset)
	if err != nil {
		return domain.TaskPage{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)
		if err != nil {
			return domain.TaskPage{}, fmt.Errorf("failed scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return domain.TaskPage{}, fmt.Errorf("failed iterating rows: %w", err)
	}

	hasMore := len(tasks) > query.Limit
	if hasMore {
		tasks = tasks[:query.Limit]
	}

	return domain.TaskPage{Tasks: toTaskDomainList(tasks), HasMore: hasMore}, nil
}

func (r *TaskSqlRepository) AddTaskToUser(userId string, task domain.Task) (domain.Task, error) {
	db := r.pgsql.connect()

	var id int
	err := db.QueryRow(`INSERT INTO task(user_id, description, created_at)
		VALUES($1, $2, $3) RETURNING id`, userId, task.Description, task.CreatedAt).Scan(&id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}

	taskEntity, err := r.GetByIDAndUserID(id, userId)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to retrieve created task: %w", err)
	}
	return toTaskDomain(Task(taskEntity)), nil
}

func (r *TaskSqlRepository) GetByIDAndUserID(id int, userId string) (domain.Task, error) {
	db := r.pgsql.connect()

	var task Task
	err := db.QueryRow(`SELECT id, user_id, description, created_at FROM task
		where id = $1 AND user_id = $2`, id, userId).Scan(&task.Id, &task.UserId, &task.Description, &task.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, domain.ErrTaskNotFound
	}

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return toTaskDomain(task), nil
}
