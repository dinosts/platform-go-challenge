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

type AudienceModel struct {
	Id                 uuid.UUID
	Gender             string
	BirthCountry       string
	AgeGroup           string
	SocialMediaHours   float64
	PurchasesLastMonth int
}

type FavouriteModel struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	AssetId     uuid.UUID
	AssetType   string
	Description string
}

type (
	UserStorage      map[uuid.UUID]UserModel
	ChartStorage     map[uuid.UUID]ChartModel
	InsightStorage   map[uuid.UUID]InsightModel
	AudienceStorage  map[uuid.UUID]AudienceModel
	FavouriteStorage map[uuid.UUID]FavouriteModel
)

type InMemoryDatabase struct {
	UserStorage      UserStorage
	ChartStorage     ChartStorage
	InsightStorage   InsightStorage
	AudienceStorage  AudienceStorage
	FavouriteStorage FavouriteStorage
}

func populateStorageForDevEnv(
	userStorage *UserStorage,
	chartStorage *ChartStorage,
	insightStorage *InsightStorage,
	audienceStorage *AudienceStorage,
	favouriteStorage *FavouriteStorage,
) {
	devUser := UserModel{Id: uuid.New(), Email: "test@test.com", Password: "pass"}
	uS := *userStorage
	uS[devUser.Id] = devUser

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
	cS := *chartStorage
	cS[chart.Id] = chart

	insight := InsightModel{
		Id:   uuid.New(),
		Text: "40% of millennials spend more than 3hours on social media daily",
	}
	iS := *insightStorage
	iS[insight.Id] = insight

	audience := AudienceModel{
		Id:                 uuid.New(),
		Gender:             "Male",
		BirthCountry:       "United Kingdom",
		AgeGroup:           "25-34",
		SocialMediaHours:   3.5,
		PurchasesLastMonth: 7,
	}
	aS := *audienceStorage
	aS[audience.Id] = audience

	fS := *favouriteStorage
	fS[uuid.New()] = FavouriteModel{
		Id:          uuid.New(),
		UserId:      devUser.Id,
		AssetId:     chart.Id,
		AssetType:   "chart",
		Description: "Main performance chart",
	}
	fS[uuid.New()] = FavouriteModel{
		Id:          uuid.New(),
		UserId:      devUser.Id,
		AssetId:     insight.Id,
		AssetType:   "insight",
		Description: "Great for Q2 presentation",
	}
	fS[uuid.New()] = FavouriteModel{
		Id:          uuid.New(),
		UserId:      devUser.Id,
		AssetId:     audience.Id,
		AssetType:   "audience",
		Description: "Target audience for campaign",
	}
}

func NewInMemoryDatabase(env string) *InMemoryDatabase {
	userStorage := UserStorage{}
	chartStorage := ChartStorage{}
	insighStorage := InsightStorage{}
	audienceStorage := AudienceStorage{}
	favouriteStorage := FavouriteStorage{}

	if env == "dev" {
		populateStorageForDevEnv(
			&userStorage,
			&chartStorage,
			&insighStorage,
			&audienceStorage,
			&favouriteStorage,
		)
	}

	return &InMemoryDatabase{
		UserStorage:    userStorage,
		ChartStorage:   chartStorage,
		InsightStorage: insighStorage,
	}
}
