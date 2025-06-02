package favourite

import (
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/insight"

	"github.com/google/uuid"
)

type AssetType string

type Favourite struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	AssetId     uuid.UUID
	AssetType   AssetType
	Description string
}

type AssetFavourites struct {
	Charts    []ChartFavourite    `json:"charts"`
	Insights  []InsightFavourite  `json:"insights"`
	Audiences []AudienceFavourite `json:"audiences"`
}

type ChartFavourite struct {
	Id          uuid.UUID   `json:"id"`
	Description string      `json:"description"`
	Info        chart.Chart `json:"info"`
}

type InsightFavourite struct {
	Id          uuid.UUID       `json:"id"`
	Description string          `json:"description"`
	Info        insight.Insight `json:"info"`
}

type AudienceFavourite struct {
	Id          uuid.UUID         `json:"id"`
	Description string            `json:"description"`
	Info        audience.Audience `json:"info"`
}

type CreateFavouriteRequestBody struct {
	AssetId     uuid.UUID `json:"assetId" validate:"required,uuid"`
	Description string    `json:"description" validate:"required"`
}

type UpdateFavouriteRequestBody struct {
	Id          uuid.UUID `json:"id" validate:"required,uuid"`
	Description string    `json:"description"`
}

type DeleteFavouriteRequestBody struct {
	Id uuid.UUID `json:"id" validate:"required,uuid"`
}
