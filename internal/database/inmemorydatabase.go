package database

import "github.com/google/uuid"

type UserModel struct {
	Id       uuid.UUID
	Email    string
	Password string
}

type ChartModel struct {
	Id         uuid.UUID
	Title      string
	XAxisTitle string
	YAxisTitle string
	Data       []map[string]float64
}

type (
	UserStorage  map[uuid.UUID]UserModel
	ChartStorage map[uuid.UUID]ChartModel
)

type InMemoryDatabase struct {
	UserStorage  UserStorage
	ChartStorage ChartStorage
}

func populateStorageForDevEnv(userStorage *UserStorage, chartStorage *ChartStorage) {
	devUser := UserModel{Id: uuid.New(), Email: "test@test.com", Password: "pass"}
	us := *userStorage
	us[devUser.Id] = devUser

	chart := ChartModel{
		Id:         uuid.New(),
		Title:      "test chart",
		XAxisTitle: "commit number",
		YAxisTitle: "lines of code",
		Data: []map[string]float64{
			{"x": 1, "y": 100},
			{"x": 2, "y": 300},
			{"x": 3, "y": 500},
		},
	}
	cs := *chartStorage
	cs[chart.Id] = chart
}

func NewInMemoryDatabase(env string) *InMemoryDatabase {
	userStorage := UserStorage{}

	if env == "dev" {
	}

	return &InMemoryDatabase{UserStorage: userStorage}
}
