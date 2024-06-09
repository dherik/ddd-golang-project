package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
)

var ErrUserNotFound = errors.New("the user was not found")

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

func (r *UserSqlRepository) GetByID(id int) (domain.User, error) {
	db, err := r.pgsql.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	var user User
	err = db.QueryRow(`SELECT id, username, email, password, created_at FROM user 
		where id = $1`, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return domain.User{}, fmt.Errorf("failed to execute get by id query: %w", err)
	}

	return toUserDomain(user), nil
}

// Add implements domain.UserRepository.
func (r *UserSqlRepository) Add(user domain.User) (domain.User, error) {
	db, err := r.pgsql.connect()
	if err != nil {
		return domain.User{}, fmt.Errorf("failed connecting to database: %w", err)
	}
	defer db.Close()

	var id int
	err = db.QueryRow(`INSERT INTO users(username, email, password, created_at)
		VALUES($1, $2, $3, $4) RETURNING id`, user.Username, user.Email, user.Password.HashPassword, user.CreatedAt).Scan(&id)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to execute insert query: %w", err)
	}

	userNew, err := r.GetByID(id)
	if err != nil {
		return domain.User{}, nil
	}
	return userNew, nil
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
		slog.Error(fmt.Sprintf("user '%s' was not found: %s", username, err.Error()))
		return domain.User{}, ErrUserNotFound
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
		Password:  domain.Password{HashPassword: userDb.Password},
		CreatedAt: userDb.CreatedAt,
	}
	return user
}
