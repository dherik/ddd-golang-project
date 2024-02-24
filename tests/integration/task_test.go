package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/dherik/ddd-golang-project/internal/domain"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TaskTestSuite struct {
	suite.Suite
	*dockertest.Pool
	*dockertest.Resource
}

func (suite *TaskTestSuite) SetupSuite() {

	if testing.Short() {
		suite.T().Skip("Skip test for postgresql repository")
	}

	dataSource, pool, resource := setup.SetupDatabase()

	suite.Pool = pool
	suite.Resource = resource

	setup.StartServer(dataSource)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExecuteAllSuites(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

func (suite *TaskTestSuite) SetupTest() {
	log.Println("##### reset database setup test")
	setup.ResetData()
	setup.LoadDML()
}

func (suite *TaskTestSuite) TearDownSuite() {
	// Teardown
	log.Println("Tear down container")

	// When you're done, kill and remove the container
	if err := suite.Pool.Purge(suite.Resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (s *TaskTestSuite) TestGetByDate() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, _ := setup.Login("admin", "some_password")

	startDate := time.Date(2024, 02, 14, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
	endDate := time.Date(2024, 02, 16, 23, 34, 58, 651387237, time.UTC).Format(time.RFC3339)

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

func (s *TaskTestSuite) TestGetByID() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, _ := setup.Login("admin", "some_password")

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/tasks/%d", 3333, 1), nil) //TODO duplicated port, get from s.port (parametrized)
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	require.JSONEq(s.T(), `{
		  "id": 1,
		  "description": "Complete project proposal",
		  "userId": "1",
		  "createdAt": "2024-02-15T10:59:01.054Z"
		}`, string(byteBody))
	response.Body.Close()
}

func (s *TaskTestSuite) TestAddTask() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, _ := setup.Login("admin", "some_password")

	payload := domain.Task{
		UserId:      "1",
		Description: "Hello, World!",
	}

	// Convert payload to JSON
	requestBody, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatalf("Error encoding JSON: %v", err)
	}

	url := fmt.Sprintf("http://localhost:%d/tasks", 3333)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody)) //TODO duplicated port, get from s.port (parametrized)
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusCreated, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal("", string(byteBody))
	response.Body.Close()

}
