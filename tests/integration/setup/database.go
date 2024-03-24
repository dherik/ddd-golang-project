package setup

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

var (
	initialized       bool
	initializedLock   sync.Mutex
	PostgresContainer *postgres.PostgresContainer
	Datasource        persistence.Datasource
	Ctx               context.Context
)

var db *sql.DB

type DatabaseIT struct {
	persistence.Datasource
}

func SetupDatabase() {

	initializedLock.Lock()
	defer initializedLock.Unlock()

	if initialized {
		slog.Debug("Database it's already initialized")
		return
	}

	slog.Info("Initializing the database...")

	datasource := persistence.Datasource{
		User:     "test_user",
		Password: "test_password",
		Name:     "test_db",
	}

	slog.Info("Starting docker")

	ctx := context.Background()
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2-alpine"),
		// postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	//FIXME it's ugly
	host, _ := postgresContainer.Host(ctx)
	datasource.Host = host
	portNat, _ := postgresContainer.MappedPort(ctx, "5432")
	port, _ := strconv.Atoi(portNat.Port())
	datasource.Port = port
	databaseUrl, _ := postgresContainer.ConnectionString(ctx, "sslmode=disable")

	slog.Info(fmt.Sprintf("Connecting to database on url: %s", databaseUrl))

	db, err = sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("Could not open SQL connection: %s", err)
	}

	slog.Info("Database for integration tests is up and running!")

	LoadDDL()

	PostgresContainer = postgresContainer
	Ctx = ctx
	Datasource = datasource

	initialized = true
}

func LoadDDL() {
	query, err := os.ReadFile("../../init_ddl.sql")
	slog.Debug("Load DDL file: " + string(query))
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

func LoadDML() {
	query, err := os.ReadFile("../../init_dml.sql")
	slog.Debug("Load DML file: " + string(query))
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}

func ResetData() {

	// Retrieve a list of all tables in the database
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		// t.Fatalf("Error querying tables: %v", err)
		panic(err)
	}
	defer rows.Close()

	// Build a comma-separated list of table names
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			// t.Fatalf("Error scanning table name: %v", err)
			panic(err)

		}
		tables = append(tables, tableName)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		// t.Fatalf("Error iterating over table names: %v", err)
		panic(err)
	}

	// If there are tables, truncate them
	if len(tables) > 0 {
		// Build the SQL query to truncate tables
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", "))

		// Execute the dynamic SQL query to truncate tables
		_, err := db.Exec(query)
		if err != nil {
			// t.Fatalf("Error resetting database: %v", err)
			panic(err)
		}

		slog.Info("Database reset successful")
	} else {
		slog.Info("No tables found in the database")
	}
}

func StopDatabase() {
	log.Println("Tear down container")

	if err := PostgresContainer.Terminate(Ctx); err != nil {
		log.Fatalf("failed to terminate postgresql container: %s", err)
	}
}
