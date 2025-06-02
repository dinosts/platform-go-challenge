package database

import (
	"errors"

	"github.com/google/uuid"
)

var ErrItemNotFound = errors.New("Not Found")

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
	// Constant UUIDs
	userId, _ := uuid.Parse("a3973a1c-a77b-4a04-a296-ddec19034419")
	chartId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")
	insightId, _ := uuid.Parse("22222222-2222-2222-2222-222222222222")
	insightId2, _ := uuid.Parse("22222222-2222-2222-2222-222222222223")
	audienceId, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")
	favChartId, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")
	favInsightId, _ := uuid.Parse("55555555-5555-5555-5555-555555555555")
	favAudienceId, _ := uuid.Parse("66666666-6666-6666-6666-666666666666")

	// User
	devUser := UserModel{
		Id:       userId,
		Email:    "test@test.com",
		Password: "pass",
	}
	(*userStorage)[devUser.Id] = devUser

	// Chart
	chart := ChartModel{
		Id:         chartId,
		Title:      "test chart",
		XAxisTitle: "commit number",
		YAxisTitle: "lines of code",
		Data: []map[string]float64{
			{"x": 1, "y": 100},
			{"x": 2, "y": 300},
			{"x": 3, "y": 500},
		},
	}
	(*chartStorage)[chart.Id] = chart

	// Insight
	insight := InsightModel{
		Id:   insightId,
		Text: "40% of millennials spend more than 3 hours on social media daily",
	}
	(*insightStorage)[insight.Id] = insight
	insight2 := InsightModel{
		Id:   insightId2,
		Text: "100% of zoomers spend more than 8 hours on watching memes",
	}
	(*insightStorage)[insight2.Id] = insight2

	// Audience
	audience := AudienceModel{
		Id:                 audienceId,
		Gender:             "Male",
		BirthCountry:       "United Kingdom",
		AgeGroup:           "25-34",
		SocialMediaHours:   3.5,
		PurchasesLastMonth: 7,
	}
	(*audienceStorage)[audience.Id] = audience

	// Favourites
	fav1 := FavouriteModel{
		Id:          favChartId,
		UserId:      devUser.Id,
		AssetId:     chart.Id,
		AssetType:   "chart",
		Description: "Main performance chart",
	}
	fav2 := FavouriteModel{
		Id:          favInsightId,
		UserId:      devUser.Id,
		AssetId:     insight.Id,
		AssetType:   "insight",
		Description: "Great for Q2 presentation",
	}
	fav3 := FavouriteModel{
		Id:          favAudienceId,
		UserId:      devUser.Id,
		AssetId:     audience.Id,
		AssetType:   "audience",
		Description: "Target audience for campaign",
	}
	(*favouriteStorage)[fav1.Id] = fav1
	(*favouriteStorage)[fav2.Id] = fav2
	(*favouriteStorage)[fav3.Id] = fav3
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
		UserStorage:      userStorage,
		ChartStorage:     chartStorage,
		InsightStorage:   insighStorage,
		AudienceStorage:  audienceStorage,
		FavouriteStorage: favouriteStorage,
	}
}

func IMStorageGetById[T any](id uuid.UUID, storage map[uuid.UUID]T) (*T, error) {
	v, found := storage[id]

	if !found {
		return nil, ErrItemNotFound
	}

	return &v, nil
}
