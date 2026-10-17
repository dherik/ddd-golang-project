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

	response, body, _ := login("non_existent_user", "some_password")

	s.Equal(http.StatusUnauthorized, response.StatusCode)
	s.assertErrorMessage(body, http.StatusUnauthorized, "invalid credentials")
}

func (s *UserTestSuite) TestLoginUnauthorizedWhenPasswordIsWrong() {

	response, body, _ := login("admin", "some_wrong_password")

	s.Equal(http.StatusUnauthorized, response.StatusCode)
	s.assertErrorMessage(body, http.StatusUnauthorized, "invalid credentials")
}

func (s *UserTestSuite) TestAddUser() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "some_user",
		Email:    "some_user@example.com",
		Password: "_1SomeValidPassword",
	}

	response := addUser(s.T(), payload, token)
	s.Equal(http.StatusCreated, response.StatusCode)

	byteBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal("", string(byteBody))
	response.Body.Close()

}

func (s *UserTestSuite) TestAddUserThenLogin() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "login_user",
		Email:    "login_user@example.com",
		Password: "_1SomeValidPassword",
	}

	response := addUser(s.T(), payload, token)
	s.Equal(http.StatusCreated, response.StatusCode)
	response.Body.Close()

	loginResponse, body, err := login(payload.Username, payload.Password)
	s.NoError(err)
	s.Equal(http.StatusOK, loginResponse.StatusCode)

	var tokenResponse setup.TokenResponse
	s.NoError(json.Unmarshal(body, &tokenResponse))
	s.NotEmpty(tokenResponse.Token)
}

func (s *UserTestSuite) TestAddUserWhenUserAlreadyExists() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "duplicate_username_user",
		Email:    "duplicate_username_user@example.com",
		Password: "_1SomeValidPassword",
	}

	first := addUser(s.T(), payload, token)
	s.Equal(http.StatusCreated, first.StatusCode)
	first.Body.Close()

	payload.Email = "some_other_user@example.com"
	response := addUser(s.T(), payload, token)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(http.StatusConflict, response.StatusCode)
	s.assertErrorMessage(body, http.StatusConflict, "a user with the same username already exists")
}

func (s *UserTestSuite) TestAddUserWhenEmailAlreadyExists() {

	token := setup.Login("admin", "some_password")

	payload := api.UserRequest{
		Username: "another_user",
		Email:    "admin@example.com",
		Password: "_1SomeValidPassword",
	}

	response := addUser(s.T(), payload, token)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	s.NoError(err)

	s.Equal(http.StatusConflict, response.StatusCode)
	s.assertErrorMessage(body, http.StatusConflict, "a user with the same email already exists")
}

func (s *UserTestSuite) assertErrorMessage(body []byte, code int, message string) {
	s.T().Helper()

	var payload struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	s.Require().NoError(json.Unmarshal(body, &payload))
	s.Equal(code, payload.Error.Code)
	s.Equal(message, payload.Error.Message)
}

func addUser(t *testing.T, payload api.UserRequest, token string) *http.Response {
	t.Helper()
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

func login(username, password string) (*http.Response, []byte, error) {
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

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp, nil, readErr
	}

	return resp, body, err
}
