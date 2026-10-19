package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TaskTestSuite struct {
	suite.Suite
}

func (s *TaskTestSuite) TestGetByDate() {

	token := setup.Login("admin", "some_password")

	startDate := time.Date(2024, 02, 14, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
	endDate := time.Date(2024, 02, 16, 23, 34, 58, 651387237, time.UTC).Format(time.RFC3339)

	url := fmt.Sprintf("http://localhost:%d/tasks?startDate=%s&endDate=%s", 3333, startDate, endDate)
	slog.Info("request url", "url", url)
	req, err := http.NewRequest(http.MethodGet, url, nil) //TODO duplicated port, get from s.port (parametrized)
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	require.JSONEq(s.T(), `{
		"items": [
			{
			  "id": 1,
			  "description": "Complete project proposal",
			  "userId": "1",
			  "createdAt": "2024-02-15T10:59:01.054Z"
			}
		],
		"hasMore": false
	  }`, string(byteBody))

	response.Body.Close()
}

func (s *TaskTestSuite) TestGetByID() {

	token := setup.Login("admin", "some_password")

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

	token := setup.Login("admin", "some_password")

	payload := api.TaskRequest{
		Description: "Hello, World!",
	}

	// Convert payload to JSON

	response := addTask(s, payload, token)
	s.Equal(http.StatusCreated, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal("", string(byteBody))
	response.Body.Close()

}

func addTask(s *TaskTestSuite, payload any, token string) *http.Response {
	s.T().Helper()
	requestBody, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatalf("Error encoding JSON: %v", err)
	}

	//TODO duplicated port, get from s.port (parametrized)
	url := fmt.Sprintf("http://localhost:%d/tasks", 3333)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	s.NoError(err)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	return response
}

func (s *TaskTestSuite) TestCannotAddTaskWhenDescriptionEmpty() {

	token := setup.Login("admin", "some_password")

	payload := api.TaskRequest{
		Description: "",
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
	s.Equal(http.StatusBadRequest, response.StatusCode)

	// byteBody, err := io.ReadAll(response.Body)
	// s.NoError(err)

	//TODO check error message
	// s.Equal("", string(byteBody))

	response.Body.Close()

}

func getTask(s *TaskTestSuite, id int, token string) *http.Response {
	s.T().Helper()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/tasks/%d", 3333, id), nil)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)
	return response
}

func (s *TaskTestSuite) TestCannotGetTaskFromAnotherUser() {

	token := setup.Login("admin", "some_password")

	response := getTask(s, 2, token)
	s.Equal(http.StatusNotFound, response.StatusCode)

	response.Body.Close()
}

func (s *TaskTestSuite) TestUserOnlySeesOwnTasks() {

	token := setup.Login("john_doe", "some_password")

	response := getTask(s, 2, token)
	s.Equal(http.StatusOK, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)
	require.JSONEq(s.T(), `{
		  "id": 2,
		  "description": "Review code for bug fixes",
		  "userId": "2",
		  "createdAt": "2024-02-16T21:23:09.066Z"
		}`, string(byteBody))
	response.Body.Close()

	response = getTask(s, 1, token)
	s.Equal(http.StatusNotFound, response.StatusCode)
	response.Body.Close()

	startDate := time.Date(2024, 02, 14, 20, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
	endDate := time.Date(2024, 02, 16, 23, 34, 58, 651387237, time.UTC).Format(time.RFC3339)
	url := fmt.Sprintf("http://localhost:%d/tasks?startDate=%s&endDate=%s", 3333, startDate, endDate)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	listResponse, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, listResponse.StatusCode)

	listBody, err := io.ReadAll(listResponse.Body)
	s.NoError(err)
	listResponse.Body.Close()

	require.JSONEq(s.T(), `{
		"items": [
			{
			  "id": 2,
			  "description": "Review code for bug fixes",
			  "userId": "2",
			  "createdAt": "2024-02-16T21:23:09.066Z"
			}
		],
		"hasMore": false
	  }`, string(listBody))
}

func (s *TaskTestSuite) TestCreateTaskIgnoresUserIdFromBody() {

	token := setup.Login("admin", "some_password")

	payload := map[string]string{
		"userId":      "2",
		"description": "Owned by admin",
	}

	response := addTask(s, payload, token)
	s.Equal(http.StatusCreated, response.StatusCode)
	response.Body.Close()

	startDate := time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)
	endDate := time.Now().UTC().Add(time.Hour).Format(time.RFC3339)
	url := fmt.Sprintf("http://localhost:%d/tasks?startDate=%s&endDate=%s", 3333, startDate, endDate)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	listResponse, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, listResponse.StatusCode)

	listBody, err := io.ReadAll(listResponse.Body)
	s.NoError(err)
	listResponse.Body.Close()

	var page api.TasksResponse
	s.NoError(json.Unmarshal(listBody, &page))

	found := false
	for _, task := range page.Items {
		if task.Description == "Owned by admin" {
			found = true
			s.Equal("1", task.UserId)
		}
	}
	s.True(found)
}
