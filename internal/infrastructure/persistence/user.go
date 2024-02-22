package persistence

import (
	"database/sql"
	"fmt"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

func NewUserRepository(postgreRepository PostgreRepository) domain.UserRepository {
	return &UserSqlRepository{pgsql: postgreRepository}
}

type UserSqlRepository struct {
	pgsql PostgreRepository
}

func (r *UserSqlRepository) FindUserByUsername(username string) (domain.User, error) {
	db, err := r.pgsql.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed connecting to database: %w", err)
	}
	defer db.Close()

	var user domain.User
	err = db.QueryRow(`SELECT id, username, email, password, created_at FROM users 
		where username = $1`, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		// TODO return specific error for user not found
		return domain.User{}, fmt.Errorf("user '%s' not found: %w", username, err)
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("failed reading rows: %w", err)
	}

	return user, nil
}
