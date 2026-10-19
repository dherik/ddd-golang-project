package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	databaseName     = "test_db"
	databaseUser     = "test_user"
	databasePassword = "test_password"
)

var (
	setupDBOnce sync.Once
	setupDBErr  error

	PostgresContainer *postgres.PostgresContainer
	Datasource        persistence.Datasource
	Ctx               context.Context
)

var db *sql.DB

// SetupDatabase starts the database container, applies the schema, and opens
// the connection the harness seeds and truncates through.
//
// It returns an error rather than exiting the process, so the suite fails with
// its own name attached and the rest of the run still reports. The one-time
// initialisation is preserved, so a failure is memoised and returned identically
// on every later call.
func SetupDatabase() error {
	setupDBOnce.Do(func() {
		setupDBErr = setupDatabase()
	})

	return setupDBErr
}

func setupDatabase() error {
	slog.Info("Initializing the database...")

	datasource := persistence.Datasource{
		User:     databaseUser,
		Password: databasePassword,
		Name:     databaseName,
	}

	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16.2-alpine",
		postgres.WithDatabase(databaseName),
		postgres.WithUsername(databaseUser),
		postgres.WithPassword(databasePassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return fmt.Errorf("failed to start the PostgreSQL container: %w", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to resolve the PostgreSQL container host: %w", err)
	}

	portNat, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return fmt.Errorf("failed to resolve the PostgreSQL container port: %w", err)
	}

	port, err := strconv.Atoi(portNat.Port())
	if err != nil {
		return fmt.Errorf("failed to parse the PostgreSQL container port %q: %w", portNat.Port(), err)
	}

	datasource.Host = host
	datasource.Port = port

	// The harness seeds the same database the application under test reads, so
	// it opens it through the production connection-string builder. That keeps
	// one description of the test database and puts the builder under test.
	connectionString := datasource.ConnectionString()
	slog.Debug("Connecting to the test database", "url", connectionString)

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("could not open the database connection: %w", err)
	}

	if err := LoadDDL(); err != nil {
		return err
	}

	slog.Info("Database for integration tests is up and running!")

	PostgresContainer = postgresContainer
	Ctx = ctx
	Datasource = datasource

	return nil
}

// LoadDDL applies the schema to the test database.
func LoadDDL() error {
	return loadSQLFile("../../init_ddl.sql")
}

// LoadDML applies the seed data to the test database.
func LoadDML() error {
	return loadSQLFile("../../init_dml.sql")
}

func loadSQLFile(path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading %s: %w", path, err)
	}

	slog.Debug("Loading SQL file", "path", path)

	if _, err := db.Exec(string(query)); err != nil {
		return fmt.Errorf("failed applying %s: %w", path, err)
	}

	return nil
}

// ResetData truncates every table so the next test starts from a known empty
// database. Truncating rather than dropping keeps the schema the harness applied.
func ResetData() error {
	// Retrieve a list of all tables in the database
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		return fmt.Errorf("failed querying tables: %w", err)
	}
	defer rows.Close()

	// Build a comma-separated list of table names
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed scanning table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed iterating over table names: %w", err)
	}

	// If there are tables, truncate them
	if len(tables) == 0 {
		slog.Info("No tables found in the database")
		return nil
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed resetting the database: %w", err)
	}

	slog.Info("Database reset successful")

	return nil
}

// StopDatabase terminates the database container.
func StopDatabase() error {
	slog.Info("Tearing down the PostgreSQL container")

	if err := PostgresContainer.Terminate(Ctx); err != nil {
		return fmt.Errorf("failed to terminate the PostgreSQL container: %w", err)
	}

	return nil
}
