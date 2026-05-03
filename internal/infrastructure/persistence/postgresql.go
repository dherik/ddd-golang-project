package persistence

import (
	"database/sql"
	"fmt"

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

type PostgreRepository struct {
	db *sql.DB
}

func NewPostgreRepository(ds Datasource) (*PostgreRepository, error) {
	db, err := sql.Open("postgres", ds.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &PostgreRepository{db: db}, nil
}

func (pg *PostgreRepository) connect() *sql.DB {
	return pg.db
}
