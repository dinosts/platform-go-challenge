package database

import "github.com/google/uuid"

type UserModel struct {
	Id       uuid.UUID
	Email    string
	Password string
}
type UserStorage map[uuid.UUID]UserModel

type InMemoryDatabase struct {
	UserStorage UserStorage
}

func NewInMemoryDatabase(env string) *InMemoryDatabase {
	userStorage := UserStorage{}

	if env == "dev" {
		devUser := UserModel{Id: uuid.New(), Email: "test@test.com", Password: "pass"}
		userStorage[devUser.Id] = devUser

	}

	return &InMemoryDatabase{UserStorage: userStorage}
}
