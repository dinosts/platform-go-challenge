package user_test

import (
	"errors"
	"platform-go-challenge/internal/domain/user"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	getByEmailFn func(email string) (*user.User, error)
}

func (m *mockUserRepository) GetByEmail(email string) (*user.User, error) {
	return m.getByEmailFn(email)
}

func TestUserService_LoginUser(t *testing.T) {
	t.Run("should return token and expiry when login is successful", func(t *testing.T) {
		// Arrange
		expectedUser := &user.User{
			Id:       uuid.New(),
			Email:    "test@example.com",
			Password: "secret123",
		}
		expectedToken := "some-token"
		expectedExpiry := time.Now().Add(time.Hour)

		mockRepo := &mockUserRepository{
			getByEmailFn: func(email string) (*user.User, error) {
				return expectedUser, nil
			},
		}
		mockTokenFn := func(claims map[string]any) (string, time.Time, error) {
			return expectedToken, expectedExpiry, nil
		}
		mockPasswordHasher := func(password string) string {
			return password
		}
		service := user.NewUserService(user.ServiceDependencies{
			UserRepository: mockRepo,
			GenerateToken:  mockTokenFn,
			PasswordHasher: mockPasswordHasher,
		})

		// Act
		token, expiresAt, err := service.LoginUser("test@example.com", "secret123")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
		assert.Equal(t, expectedExpiry, expiresAt)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		// Arrange
		mockRepo := &mockUserRepository{
			getByEmailFn: func(email string) (*user.User, error) {
				return nil, user.ErrUserNotFound
			},
		}
		service := user.NewUserService(user.ServiceDependencies{
			UserRepository: mockRepo,
			GenerateToken:  nil,
			PasswordHasher: nil,
		})

		// Act
		token, expiresAt, err := service.LoginUser("unknown@example.com", "whatever")

		// Assert
		assert.ErrorIs(t, err, user.ErrLoginFailed)
		assert.Empty(t, token)
		assert.True(t, expiresAt.IsZero())
	})

	t.Run("should return error when password does not match", func(t *testing.T) {
		// Arrange
		mockRepo := &mockUserRepository{
			getByEmailFn: func(email string) (*user.User, error) {
				return &user.User{
					Id:       uuid.New(),
					Email:    email,
					Password: "correctpass",
				}, nil
			},
		}
		service := user.NewUserService(user.ServiceDependencies{
			UserRepository: mockRepo,
			GenerateToken:  nil,
			PasswordHasher: func(password string) string { return password },
		})

		// Act
		token, expiresAt, err := service.LoginUser("test@example.com", "wrongpass")

		// Assert
		assert.ErrorIs(t, err, user.ErrLoginFailed)
		assert.Empty(t, token)
		assert.True(t, expiresAt.IsZero())
	})

	t.Run("should return error when token generation fails", func(t *testing.T) {
		// Arrange
		mockRepo := &mockUserRepository{
			getByEmailFn: func(email string) (*user.User, error) {
				return &user.User{
					Id:       uuid.New(),
					Email:    email,
					Password: "pass",
				}, nil
			},
		}
		mockTokenFn := func(claims map[string]any) (string, time.Time, error) {
			return "", time.Time{}, errors.New("token failure")
		}
		service := user.NewUserService(user.ServiceDependencies{
			UserRepository: mockRepo,
			GenerateToken:  mockTokenFn,
			PasswordHasher: func(password string) string { return password },
		})

		// Act
		token, expiresAt, err := service.LoginUser("test@example.com", "pass")

		// Assert
		assert.ErrorIs(t, err, user.ErrTokenGenerationFailed)
		assert.Empty(t, token)
		assert.True(t, expiresAt.IsZero())
	})
}
