package persistence

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

type User struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

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

	var user User
	err = db.QueryRow(`SELECT id, username, email, password, created_at FROM users 
		where username = $1`, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		// TODO return specific error for user not found
		return domain.User{}, fmt.Errorf("user '%s' not found: %w", username, err)
	}

	if err != nil {
		return domain.User{}, fmt.Errorf("failed reading rows: %w", err)
	}

	return toUserDomain(user), nil
}

func toUserDomain(userDb User) domain.User {
	user := domain.User{
		Id:        userDb.Id,
		Username:  userDb.Username,
		Email:     userDb.Email,
		Password:  userDb.Password,
		CreatedAt: userDb.CreatedAt,
	}
	return user
}
