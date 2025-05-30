package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

func NewJWTAuth(secret string) *jwtauth.JWTAuth {
	jwtAuth := jwtauth.New("HS256", []byte(secret), nil)

	debugToken, _ := GenerateJWToken(jwtAuth, map[string]any{})

	fmt.Print(debugToken)

	return jwtAuth
}

func AuthenticatorMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				message := fmt.Sprintf("Could not Authorize: %s", err.Error())
				RespondWithError(w, http.StatusUnauthorized, message)
				return
			}

			if token == nil {
				message := "Authorization token not found"
				RespondWithError(w, http.StatusUnauthorized, message)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(handler)
	}
}

func VerifierMiddleware(jwtAuth *jwtauth.JWTAuth) func(next http.Handler) http.Handler {
	return jwtauth.Verifier(jwtAuth)
}

func GenerateJWToken(tokenAuth *jwtauth.JWTAuth, extraTokenInfo map[string]any) (string, error) {
	tokenInfo := map[string]interface{}{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	for key, value := range extraTokenInfo {
		tokenInfo[key] = value
	}

	_, tokenString, err := tokenAuth.Encode(tokenInfo)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return tokenString, nil
}
