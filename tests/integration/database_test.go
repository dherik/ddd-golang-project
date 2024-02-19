package integration

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func setupDatabase(suite *ExampleTestSuite) persistence.Datasource {

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
		Tag:        "16.2",
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

	loadTestData()

	//Run tests
	// code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	// if err := pool.Purge(resource); err != nil {
	// 	log.Fatalf("Could not purge resource: %s", err)
	// }

	// os.Exit(code)

	suite.Pool = *pool
	suite.Resource = resource

	return datasource
}

func loadTestData() {
	query, err := os.ReadFile("../../init.sql")
	log.Printf("Load file: " + string(query))
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
}
