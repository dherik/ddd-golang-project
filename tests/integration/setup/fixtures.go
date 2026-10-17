package setup

import "time"

// InsertTaskWithCreatedAt inserts a task with an explicit creation time so that
// ordering tests can control placement, including timestamp ties.
func InsertTaskWithCreatedAt(userId int, description string, createdAt time.Time) int {
	var id int
	err := db.QueryRow(
		`INSERT INTO task (user_id, description, created_at) VALUES ($1, $2, $3) RETURNING id`,
		userId, description, createdAt,
	).Scan(&id)
	if err != nil {
		panic(err)
	}
	return id
}
