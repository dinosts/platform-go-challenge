package user_test

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDBUserRepository_GetByEmail(t *testing.T) {
	t.Run("should return user when email exists in database", func(t *testing.T) {
		// Arrange
		expectedUser := database.UserModel{
			Id:       uuid.New(),
			Email:    "test@example.com",
			Password: "secret",
		}
		db := &database.InMemoryDatabase{
			UserStorage: database.UserStorage{
				expectedUser.Id: expectedUser,
			},
		}
		repo := user.NewInMemoryDBUserRepository(db)

		// Act
		result, err := repo.GetByEmail("test@example.com")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.Id, result.Id)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.Password, result.Password)
	})

	t.Run("should return error when email does not exist in database", func(t *testing.T) {
		// Arrange
		db := &database.InMemoryDatabase{
			UserStorage: database.UserStorage{},
		}
		repo := user.NewInMemoryDBUserRepository(db)

		// Act
		result, err := repo.GetByEmail("notfound@example.com")

		// Assert
		assert.Nil(t, result)
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})
}
