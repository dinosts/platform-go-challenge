package e2e

import (
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/server"
	"testing"
)

func TestGetHealth(t *testing.T) {
	// Arrange
	router := server.SetupRouter()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	expected_status := http.StatusOK
	expected_body := "Healthy!"

	// Act
	router.ServeHTTP(res, req)

	// Assert
	if res.Result().StatusCode != expected_status {
		t.Fatalf("expected status %d, got %d", expected_status, res.Result().StatusCode)
	}

	body := string(res.Body.Bytes())
	if body != expected_body {
		t.Fatalf("expected body %s, got %s", expected_body, body)
	}
}
