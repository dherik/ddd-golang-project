package setup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func Login(username, password string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add fields to the form
	_ = writer.WriteField("username", username)
	_ = writer.WriteField("password", password)

	err := writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/login", body) //FIXME port as parameter
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		var tokenResponse TokenResponse
		err := json.Unmarshal(responseBody, &tokenResponse)
		if err != nil {
			return "", err
		}
		return tokenResponse.Token, nil
	}

	return "", fmt.Errorf("login failed. Status code: %d, Response: %s", resp.StatusCode, string(responseBody))
}
