package integration

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

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
