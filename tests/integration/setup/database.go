package setup

import (
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
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	initialized     bool
	initializedLock sync.Mutex
	Pool            *dockertest.Pool
	Resource        *dockertest.Resource
	Datasource      persistence.Datasource
)

var db *sql.DB

type DatabaseIT struct {
	persistence.Datasource
	*dockertest.Resource
}

func SetupDatabase() {

	initializedLock.Lock()
	defer initializedLock.Unlock()

	if initialized {
		slog.Debug("Database it's already initialized")
		return
	}

	fmt.Println("Initializing the database...") //FIXME slog

	datasource := persistence.Datasource{
		User:     "test_user",
		Password: "test_password",
		Name:     "test_db",
	}

	slog.Info("starting docker")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2", // TODO use the same from docker-compose?
		Env: []string{
			"POSTGRES_PASSWORD=test_password",
			"POSTGRES_USER=test_user",
			"POSTGRES_DB=test_db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	hostAndPortArray := strings.Split(hostAndPort, ":")
	datasource.Host = hostAndPortArray[0]
	port, _ := strconv.Atoi(hostAndPortArray[1])
	datasource.Port = port
	databaseUrl := datasource.ConnectionString()

	slog.Info(fmt.Sprintf("Connecting to database on url: %s", databaseUrl))

	// code := m.Run()
	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	slog.Info("Database for integration tests is up and running!")

	LoadDDL()

	Pool = pool
	Resource = resource
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

	if err := Pool.Purge(Resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
