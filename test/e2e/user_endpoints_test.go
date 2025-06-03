package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"platform-go-challenge/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	// Arrange
	server, _ := test.StartServer()
	defer server.Close()

	client := server.Client()

	requestBody := map[string]interface{}{
		"email":    "test@test.com",
		"password": "pass",
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, server.URL+"/v1/user/login", bytes.NewReader(bodyBytes))

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	data := result["data"].(map[string]any)
	assert.NotEmpty(t, data["token"])
	assert.NotEmpty(t, data["expires_at"])
}
