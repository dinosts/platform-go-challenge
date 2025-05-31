package user

import (
	"platform-go-challenge/internal/database"
)

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
}

type InMemoryDBUserRepository struct {
	DB *database.InMemoryDatabase
}

func InMemoryUserModelToDTO(userModel database.UserModel) User {
	return User{
		Id:       userModel.Id,
		Email:    userModel.Email,
		Password: userModel.Password,
	}
}

func (repo *InMemoryDBUserRepository) GetUserByEmail(email string) (*User, error) {
	for _, user := range repo.DB.UserStorage {
		if user.Email == email {
			userDTO := InMemoryUserModelToDTO(user)
			return &userDTO, nil
		}
	}

	return nil, ErrUserNotFound
}
