package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func NewJWTAuth(secret string) *jwtauth.JWTAuth {
	jwtAuth := jwtauth.New("HS256", []byte(secret), nil)

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

func NewJWToken(tokenAuth *jwtauth.JWTAuth, extraTokenInfo map[string]any) (string, time.Time, error) {
	tokenInfo := map[string]interface{}{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	for key, value := range extraTokenInfo {
		tokenInfo[key] = value
	}

	token, tokenString, err := tokenAuth.Encode(tokenInfo)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return tokenString, token.Expiration(), nil
}

func NewJWTokenIssuer(tokenAuth *jwtauth.JWTAuth) func(extraTokenInfo map[string]any) (string, time.Time, error) {
	return func(extraTokenInfo map[string]any) (string, time.Time, error) {
		return NewJWToken(tokenAuth, extraTokenInfo)
	}
}

func GetJWTokenSub(r *http.Request) (string, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("sub claim not found or not a string")
	}

	return sub, nil
}

func GetUserIdFromAuthToken(r *http.Request) (uuid.UUID, error) {
	tokenString, err := GetJWTokenSub(r)
	if err != nil {
		return uuid.Nil, err
	}

	tokenUUID, err := uuid.Parse(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	return tokenUUID, nil
}
