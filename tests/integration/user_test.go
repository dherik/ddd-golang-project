package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/tests/integration/setup"
	"github.com/stretchr/testify/suite"
)

//TODO try to access protect endpoint not authenticated test
//TODO try to access protect endpoint authenticated test
//TODO create user and login test

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) TestLoginSuccess() {

	token := setup.Login("admin", "some_password")

	if token == "" {
		s.T().Fatalf("token should not be empty")
	}
}

func (s *UserTestSuite) TestLoginUnauthorizedWhenUserNotFound() {

	response, _ := login("non_existent_user", "some_password")
	s.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *UserTestSuite) TestLoginUnauthorizedWhenPasswordIsWrong() {

	response, _ := login("admin", "some_wrong_password")
	s.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *UserTestSuite) TestAddUser() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "some_user",
		Email:    "some_user@example.com",
		Password: "some_user_password",
	}

	response := addUser(s.T(), payload, token)
	s.Equal(http.StatusCreated, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal("", string(byteBody))
	response.Body.Close()

}

func (s *UserTestSuite) TestAddUserWhenUserAlreadyExists() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "some_user",
		Email:    "some_user@example.com",
		Password: "some_user_password",
	}

	addUser(s.T(), payload, token)

	// add same user again
	response := addUser(s.T(), payload, token)

	s.Equal(http.StatusConflict, response.StatusCode)
}

// func (s *UserTestSuite) TestCannotAccessTaskFromAnotherUser() {

// 	token := setup.Login("admin", "some_password")

// 	task1 := api.TaskRequest{
// 		UserId:      "1",
// 		Description: "Hello, World!",
// 	}
// 	addTask(s.T(), task1, token)

// }

func addUser(t *testing.T, payload api.UserRequest, token string) *http.Response {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Error encoding JSON: %v", err)
	}

	//TODO duplicated port, get from s.port (parametrized)
	url := fmt.Sprintf("http://localhost:%d/users", 3333)
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error adding test user: %v", err)
	}
	return response
}

func login(username, password string) (*http.Response, error) {
	creds := api.Credentials{
		Username: username,
		Password: password,
	}

	// Convert Credentials struct to JSON
	jsonData, err := json.Marshal(creds)
	if err != nil {
		log.Fatalf("failed to marshal JSON: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/login", bytes.NewBuffer(jsonData)) //FIXME port as parameter
	if err != nil {
		log.Fatalf("failed to create the login request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed doing login request: %s", err.Error())
	}
	defer resp.Body.Close()

	return resp, err

	// if resp.StatusCode == http.StatusOK {
	// 	var tokenResponse TokenResponse
	// 	err := json.Unmarshal(responseBody, &tokenResponse)
	// 	if err != nil {
	// 		log.Fatalf("failed unmarhalling login response body: %s", err.Error())
	// 	}
	// 	return tokenResponse.Token
	// }

	// log.Fatalf("login failed. Status code: %d, Response: %s", resp.StatusCode, string(responseBody))

	// return ""
}
