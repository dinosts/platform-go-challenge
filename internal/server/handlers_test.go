package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/server"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	expectedStatus := http.StatusOK
	expectedResponseBody := utils.MessageResponse{Message: "Healthy!"}

	handler := http.HandlerFunc(server.GetHealth)

	// Act
	handler.ServeHTTP(res, req)

	// Assert
	var parsedBody utils.MessageResponse
	err := json.NewDecoder(res.Body).Decode(&parsedBody)
	assert.NoError(t, err)

	assert.Equal(t, expectedStatus, res.Code)
	assert.Equal(t, expectedResponseBody, parsedBody)
}
