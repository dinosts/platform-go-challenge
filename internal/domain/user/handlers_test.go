package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/domain/user"
	"platform-go-challenge/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock UserService
type mockUserService struct {
	loginFn func(email, password string) (string, time.Time, error)
}

func (m *mockUserService) LoginUser(email, password string) (string, time.Time, error) {
	return m.loginFn(email, password)
}

func TestUserLoginHandler(t *testing.T) {
	t.Run(
		"test should return token and expires at",
		func(t *testing.T) {
			// Arrange
			expectedExpiresAt := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
			expectedToken := "valid token"
			requestBody := map[string]string{
				"email":    "test@exapmle.com",
				"password": "secret123",
			}
			loginFn := func(email, password string) (string, time.Time, error) {
				return expectedToken, expectedExpiresAt, nil
			}

			var req *http.Request
			body, _ := json.Marshal(requestBody)
			req = httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			handler := user.UserLoginHandler(user.UserLoginDependencies{
				UserService: &mockUserService{loginFn: loginFn},
			})

			expectedStatus := http.StatusOK
			expectedResponseBody := utils.DataResponse[user.UserLoginResponseBody]{
				Data: user.UserLoginResponseBody{
					Token:     expectedToken,
					ExpiresAt: expectedExpiresAt,
				},
			}

			// Act
			handler.ServeHTTP(res, req)

			// Assert
			var parsedBody utils.DataResponse[user.UserLoginResponseBody]
			err := json.NewDecoder(res.Body).Decode(&parsedBody)
			assert.NoError(t, err)

			assert.Equal(t, expectedStatus, res.Code)
			assert.Equal(t, expectedResponseBody, parsedBody)
		},
	)

	errorResponseTests := []struct {
		name                 string
		requestBody          map[string]string
		loginFn              func(email, password string) (string, time.Time, error)
		expectedStatus       int
		expectedResponseBody utils.ErrorResponse
	}{
		{
			name: "should return Unauthorized login when failed to login",
			requestBody: map[string]string{
				"email":    "wrong@example.com",
				"password": "wrongpass",
			},
			loginFn: func(email, password string) (string, time.Time, error) {
				return "", time.Time{}, user.ErrLoginFailed
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedResponseBody: utils.ErrorResponse{Error: "Invalid email or password"},
		},
		{
			name: "should return invalid request body when request body wrong format",
			requestBody: map[string]string{
				"email":    "not an email",
				"password": "wrongpass",
			},
			loginFn:              nil,
			expectedStatus:       http.StatusBadRequest,
			expectedResponseBody: utils.ErrorResponse{Error: "Body Validation Failed, Key: 'UserLoginRequestBody.Email' Error:Field validation for 'Email' failed on the 'email' tag"},
		},
		{
			name:           "should return invalid request body when no request body",
			requestBody:    nil,
			loginFn:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponseBody: utils.ErrorResponse{
				Error: "Body Validation Failed, Key: 'UserLoginRequestBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\n" +
					"Key: 'UserLoginRequestBody.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
	}
	for _, testData := range errorResponseTests {
		t.Run(testData.name, func(t *testing.T) {
			// Arrange
			var req *http.Request
			body, _ := json.Marshal(testData.requestBody)
			req = httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			handler := user.UserLoginHandler(user.UserLoginDependencies{
				UserService: &mockUserService{loginFn: testData.loginFn},
			})

			// Act
			handler.ServeHTTP(res, req)

			// Assert
			var parsedBody utils.ErrorResponse
			err := json.NewDecoder(res.Body).Decode(&parsedBody)
			assert.NoError(t, err)

			assert.Equal(t, testData.expectedStatus, res.Code)
			assert.Equal(t, testData.expectedResponseBody, parsedBody)
		})
	}
}
