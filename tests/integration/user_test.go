package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

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

	requestBody, err := json.Marshal(payload)
	if err != nil {
		s.T().Fatalf("Error encoding JSON: %v", err)
	}

	url := fmt.Sprintf("http://localhost:%d/users", 3333)                  //FIXME port
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

// func (s *UserTestSuite) TestAddUserWhenUserAlreadyExists() {

// }

func login(username, password string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add fields to the form
	_ = writer.WriteField("username", username)
	_ = writer.WriteField("password", password)

	err := writer.Close()
	if err != nil {
		log.Fatalf("failed to close writer: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/login", body) //FIXME port as parameter
	if err != nil {
		log.Fatalf("failed to create the login request: %s", err.Error())
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
