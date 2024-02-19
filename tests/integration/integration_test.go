package integration

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var db *sql.DB

type ExampleTestSuite struct {
	suite.Suite
	dockertest.Pool
	*dockertest.Resource
}

func (suite *ExampleTestSuite) SetupSuite() {

	if testing.Short() {
		suite.T().Skip("Skip test for mongodb repository")
	}

	dataSource := setupDatabase(suite)
	startServer(dataSource)
}

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

func startServer(dataSource persistence.Datasource) {
	serverReady := make(chan bool)

	server := app.Server{
		Datasource:  dataSource,
		ServerReady: serverReady,
	}

	go server.Start()
	waitServiceStart("http://localhost:3333", 20, 100*time.Millisecond)
}

func waitServiceStart(url string, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		// Make a GET request to the /api/health endpoint
		response, err := http.Get(url + "/api/health")
		if err == nil && response.StatusCode == http.StatusOK {
			// The service is healthy, proceed with your logic
			slog.Info("Service is healthy!")
			return nil
		}

		// Print an error message (optional)
		slog.Info(fmt.Sprintf("Attempt %d: Service is not healthy yet. Retrying in %s...\n", i+1, retryInterval))

		// Wait for the specified interval before retrying
		time.Sleep(retryInterval)
	}

	// Return an error if the service does not become healthy within the specified retries
	return fmt.Errorf("Service did not become healthy after %d attempts", maxRetries)
}

func (suite *ExampleTestSuite) TearDownSuite() {
	// Teardown
	log.Println("Tear down container")

	// When you're done, kill and remove the container
	if err := suite.Pool.Purge(suite.Resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (s *ExampleTestSuite) TestGetByDate() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, _ := login()

	startDate := time.Now().AddDate(0, 0, -5).UTC().Format(time.RFC3339)
	endDate := time.Now().AddDate(0, 0, 5).UTC().Format(time.RFC3339)

	url := fmt.Sprintf("http://localhost:%d/tasks?startDate=%s&endDate=%s", 3333, startDate, endDate)
	slog.Info(fmt.Sprintf("Url: %s", url))
	req, err := http.NewRequest(http.MethodGet, url, nil) //TODO duplicated port, get from s.port (parametrized)
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	require.JSONEq(s.T(), `[
		{
		  "id": 1,
		  "description": "Complete project proposal",
		  "userId": "1",
		  "createdAt": "2024-02-15T10:59:01.054Z"
		},
		{
		  "id": 2,
		  "description": "Review code for bug fixes",
		  "userId": "2",
		  "createdAt": "2024-02-16T21:23:09.066Z"
		}
	  ]`, string(byteBody))

	response.Body.Close()
}

func (s *ExampleTestSuite) TestGetByID() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, _ := login()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/tasks/%d", 3333, 1), nil) //TODO duplicated port, get from s.port (parametrized)
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	require.JSONEq(s.T(), `[
		{
		  "id": 1,
		  "description": "Complete project proposal",
		  "userId": "1",
		  "createdAt": "2024-02-15T10:59:01.054Z"
		}
	  ]`, string(byteBody))
	response.Body.Close()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
