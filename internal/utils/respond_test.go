package utils_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	// Arrange
	recorder := httptest.NewRecorder()

	// Act
	utils.RespondWithError(recorder, http.StatusBadRequest, "Bad Request")

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var body utils.ErrorResponse
	err := json.NewDecoder(recorder.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "Bad Request", body.Error)
}

func TestRespondWithMessage(t *testing.T) {
	// Arrange
	recorder := httptest.NewRecorder()

	// Act
	utils.RespondWithMessage(recorder, http.StatusOK, "Success")

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var body utils.MessageResponse
	err := json.NewDecoder(recorder.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, "Success", body.Message)
}

func TestRespondWithData(t *testing.T) {
	// Arrange
	recorder := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	// Act
	utils.RespondWithData(recorder, http.StatusCreated, data)

	// Assert
	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var body utils.DataResponse[any]
	err := json.NewDecoder(recorder.Body).Decode(&body)
	assert.NoError(t, err)

	decodedData, ok := body.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", decodedData["key"])
}
