package e2e

import (
	"encoding/json"
	"net/http"
	"platform-go-challenge/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	// Arrange
	server, _ := test.StartServer()
	defer server.Close()

	client := server.Client()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Healthy!", result["message"])
}
