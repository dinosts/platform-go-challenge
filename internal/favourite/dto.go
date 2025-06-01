package favourite

import "github.com/google/uuid"

type AssetType string

type Favourite struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	AssetId     uuid.UUID
	AssetType   AssetType
	Description string
}

type GroupedFavourites struct {
	Insights []Favourite
	Charts   []Favourite
	Audience []Favourite
}

type PaginatedFavourites struct {
	Items    []Favourite
	Page     int
	PageSize int
	MaxPage  int
}

type ChartFavourite {
	Id          uuid.UUID
	Description string
}
