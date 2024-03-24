package setup

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/dherik/ddd-golang-project/internal/app/api"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func Login(username, password string) string {

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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed reading login body request: %s", err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var tokenResponse TokenResponse
		err := json.Unmarshal(responseBody, &tokenResponse)
		if err != nil {
			log.Fatalf("failed unmarhalling login response body: %s", err.Error())
		}
		return tokenResponse.Token
	}

	log.Fatalf("login failed. Status code: %d, Response: %s", resp.StatusCode, string(responseBody))

	return ""
}
