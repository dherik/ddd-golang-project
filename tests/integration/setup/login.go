package setup

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func Login(username, password string) string {
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
