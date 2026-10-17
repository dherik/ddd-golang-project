package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/lib/pq"
)

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUserRepository(postgreRepository *PostgreRepository) domain.UserRepository {
	return &UserSqlRepository{pgsql: postgreRepository}
}

type UserSqlRepository struct {
	pgsql *PostgreRepository
}

// Add implements domain.UserRepository.
func (r *UserSqlRepository) Add(user domain.User) error {
	db := r.pgsql.connect()

	_, err := db.Exec(`INSERT INTO users(username, email, password, created_at)
		VALUES($1, $2, $3, $4)`, user.Username, user.Email, user.Password.Hash(), user.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "unique_username":
				return domain.ErrUsernameAlreadyExists
			case "unique_email":
				return domain.ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (r *UserSqlRepository) FindUserByUsername(username string) (domain.User, error) {
	db := r.pgsql.connect()

	var user User
	err := db.QueryRow(`SELECT id, username, email, password, created_at FROM users
		where username = $1`, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrUserNotFound
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("failed reading rows: %w", err)
	}

	return toUserDomain(user), nil
}

func toUserDomain(userDb User) domain.User {
	return domain.User{
		Id:        userDb.Id,
		Username:  userDb.Username,
		Email:     userDb.Email,
		Password:  domain.RehydratePassword(userDb.Password),
		CreatedAt: userDb.CreatedAt,
	}
}
