package database

import (
	"errors"

	"github.com/google/uuid"
)

var IMErrItemNotFound = errors.New("Not Found")

type IMUserModel struct {
	Id       uuid.UUID
	Email    string
	Password string
}

type IMInsightModel struct {
	Id   uuid.UUID
	Text string
}

type IMChartModel struct {
	Id         uuid.UUID
	Title      string
	XAxisTitle string
	YAxisTitle string
	Data       []map[string]float64
}

type IMAudienceModel struct {
	Id                 uuid.UUID
	Gender             string
	BirthCountry       string
	AgeGroup           string
	SocialMediaHours   float64
	PurchasesLastMonth int
}

type IMFavouriteModel struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	AssetId     uuid.UUID
	AssetType   string
	Description string
}

type (
	UserStorage      map[uuid.UUID]IMUserModel
	ChartStorage     map[uuid.UUID]IMChartModel
	InsightStorage   map[uuid.UUID]IMInsightModel
	AudienceStorage  map[uuid.UUID]IMAudienceModel
	FavouriteStorage map[uuid.UUID]IMFavouriteModel
)

type IMDatabase struct {
	UserStorage      UserStorage
	ChartStorage     ChartStorage
	InsightStorage   InsightStorage
	AudienceStorage  AudienceStorage
	FavouriteStorage FavouriteStorage
}

func NewIMDatabase() *IMDatabase {
	userStorage := UserStorage{}
	chartStorage := ChartStorage{}
	insighStorage := InsightStorage{}
	audienceStorage := AudienceStorage{}
	favouriteStorage := FavouriteStorage{}

	return &IMDatabase{
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
		return nil, IMErrItemNotFound
	}

	return &v, nil
}
