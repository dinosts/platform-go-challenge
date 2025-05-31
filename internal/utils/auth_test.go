package utils_test

import (
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/utils"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWToken(t *testing.T) {
	// Arrange
	tokenAuth := utils.NewJWTAuth("secret")

	// Act
	token, expiresAt, err := utils.NewJWToken(tokenAuth, map[string]any{"user_id": 123})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, expiresAt)
}

func TestVerifierAndAuthenticatorMiddleware_ValidToken(t *testing.T) {
	// Arrange
	jwtAuth := jwtauth.New("HS256", []byte("secret"), nil)

	handler := setupHandler(jwtAuth)
	token, _, _ := utils.NewJWToken(jwtAuth, map[string]any{"user_id": 1})

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+token)

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Authorized", strings.TrimSpace(w.Body.String()))
}

func TestAuthenticatorMiddleware_InvalidToken(t *testing.T) {
	handler := setupHandler(nil)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer invalid.token.here")

	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticatorMiddleware_MissingToken(t *testing.T) {
	// Arrange
	handler := setupHandler(nil)
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Act
	handler.ServeHTTP(w, r)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// Helpers
func setupHandler(jwtAuth *jwtauth.JWTAuth) http.Handler {
	if jwtAuth == nil {
		jwtAuth = jwtauth.New("HS256", []byte("secret"), nil)
	}

	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Authorized"))
	})

	authHandler := utils.AuthenticatorMiddleware()(baseHandler)

	return utils.VerifierMiddleware(jwtAuth)(authHandler)
}
