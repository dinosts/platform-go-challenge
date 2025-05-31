package database

import "github.com/google/uuid"

type UserModel struct {
	Id       uuid.UUID
	Email    string
	Password string
}

type InsightModel struct {
	Id   uuid.UUID
	Text string
}

type ChartModel struct {
	Id         uuid.UUID
	Title      string
	XAxisTitle string
	YAxisTitle string
	Data       []map[string]float64
}

type (
	UserStorage    map[uuid.UUID]UserModel
	ChartStorage   map[uuid.UUID]ChartModel
	InsightStorage map[uuid.UUID]InsightModel
)

type InMemoryDatabase struct {
	UserStorage    UserStorage
	ChartStorage   ChartStorage
	InsightStorage InsightStorage
}

func populateStorageForDevEnv(
	userStorage *UserStorage,
	chartStorage *ChartStorage,
	insightStorage *InsightStorage,
) {
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

	insight := InsightModel{
		Id:   uuid.New(),
		Text: "40% of millennials spend more than 3hours on social media daily",
	}
	is := *insightStorage
	is[insight.Id] = insight
}

func NewInMemoryDatabase(env string) *InMemoryDatabase {
	userStorage := UserStorage{}
	chartStorage := ChartStorage{}
	insighStorage := InsightStorage{}

	if env == "dev" {
		populateStorageForDevEnv(
			&userStorage,
			&chartStorage,
			&insighStorage,
		)
	}

	return &InMemoryDatabase{
		UserStorage:    userStorage,
		ChartStorage:   chartStorage,
		InsightStorage: insighStorage,
	}
}
