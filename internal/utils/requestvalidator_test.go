package utils_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func TestBodyValidator(t *testing.T) {
	t.Run("should parse and inject valid body into context", func(t *testing.T) {
		// Arrange
		body := map[string]string{
			"email":    "valid@example.com",
			"password": "supersecret",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		var capturedParsed DummyRequest

		// Dummy handler that reads from context
		handler := utils.BodyValidator[DummyRequest](
			func(w http.ResponseWriter, r *http.Request) {
				val, ok := utils.GetParsedBody[DummyRequest](r)
				assert.True(t, ok)
				capturedParsed = val
			},
		)

		// Act
		handler.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, body["email"], capturedParsed.Email)
		assert.Equal(t, body["password"], capturedParsed.Password)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		// Arrange
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler := utils.BodyValidator[DummyRequest](func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not call next handler on invalid JSON")
		})

		// Act
		handler.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should return 400 for validation errors", func(t *testing.T) {
		// Arrange
		body := map[string]string{
			"email":    "invalid-email",
			"password": "123",
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()

		handler := utils.BodyValidator[DummyRequest](func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("should not call next handler on validation error")
		})

		// Act
		handler.ServeHTTP(res, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}

func TestGetParsedBody(t *testing.T) {
	t.Run("should return false when context has no parsed body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		val, ok := utils.GetParsedBody[DummyRequest](req)
		assert.False(t, ok)
		assert.Equal(t, DummyRequest{}, val)
	})
}
