package user

import (
	"platform-go-challenge/internal/database"
)

type UserRepository interface {
	GetByEmail(email string) (*User, error)
}

type inMemoryDBUserRepository struct {
	DB *database.IMDatabase
}

func NewInMemoryDBUserRepository(db *database.IMDatabase) inMemoryDBUserRepository {
	return inMemoryDBUserRepository{
		DB: db,
	}
}

func InMemoryDBUserModelToDTO(userModel database.IMUserModel) User {
	return User{
		Id:       userModel.Id,
		Email:    userModel.Email,
		Password: userModel.Password,
	}
}

func (repo *inMemoryDBUserRepository) GetByEmail(email string) (*User, error) {
	for _, user := range repo.DB.UserStorage {
		if user.Email == email {
			userDTO := InMemoryDBUserModelToDTO(user)
			return &userDTO, nil
		}
	}

	return nil, ErrUserNotFound
}
