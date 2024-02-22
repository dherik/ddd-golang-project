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
	DataSource Datasource
}

func NewPostgreRepository(ds Datasource) PostgreRepository {
	return PostgreRepository{
		DataSource: ds,
	}
}

func (pg *PostgreRepository) connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", pg.DataSource.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed connecting to database: %w", err)
	}

	return db, nil
}
