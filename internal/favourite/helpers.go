package favourite

import (
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/insight"

	"github.com/google/uuid"
)

func ExtractAssetTypeIds(assetType AssetType, favourites []Favourite) uuid.UUIDs {
	result := uuid.UUIDs{}

	for _, favourite := range favourites {
		if assetType == favourite.AssetType {
			result = append(result, favourite.AssetId)
		}
	}

	return result
}

func BuildAssetFavourites(
	favourites []Favourite,
	charts []chart.Chart,
	insights []insight.Insight,
	audiences []audience.Audience,
) (*AssetFavourites, error) {
	chartsFavourites, err := buildChartFavourites(charts, favourites)
	if err != nil {
		return nil, err
	}

	insightFavourites, err := buildInsightFavourites(insights, favourites)
	if err != nil {
		return nil, err
	}

	audienceFavourites, err := buildAudienceFavourite(audiences, favourites)
	if err != nil {
		return nil, err
	}

	return &AssetFavourites{
		Charts:    *chartsFavourites,
		Insights:  *insightFavourites,
		Audiences: *audienceFavourites,
	}, nil
}

func findFavouriteInSliceByAssetId(id uuid.UUID, favourites []Favourite) *Favourite {
	for _, favourite := range favourites {
		if favourite.AssetId == id {
			return &favourite
		}
	}

	return nil
}

func buildChartFavourites(charts []chart.Chart, favourites []Favourite) (*[]ChartFavourite, error) {
	result := []ChartFavourite{}

	for _, chart := range charts {
		favourite := findFavouriteInSliceByAssetId(chart.Id, favourites)

		if favourite == nil {
			return nil, ErrCouldNotFindFavouriteForAsset
		}

		result = append(
			result,
			ChartFavourite{
				Id:          favourite.Id,
				Description: favourite.Description,
				Info:        chart,
			},
		)
	}

	return &result, nil
}

func buildInsightFavourites(insights []insight.Insight, favourites []Favourite) (*[]InsightFavourite, error) {
	result := []InsightFavourite{}

	for _, insight := range insights {
		favourite := findFavouriteInSliceByAssetId(insight.Id, favourites)

		if favourite == nil {
			return nil, ErrCouldNotFindFavouriteForAsset
		}

		result = append(
			result,
			InsightFavourite{
				Id:          favourite.Id,
				Description: favourite.Description,
				Info:        insight,
			},
		)
	}

	return &result, nil
}

func buildAudienceFavourite(audiences []audience.Audience, favourites []Favourite) (*[]AudienceFavourite, error) {
	result := []AudienceFavourite{}

	for _, audience := range audiences {
		favourite := findFavouriteInSliceByAssetId(audience.Id, favourites)

		if favourite == nil {
			return nil, ErrCouldNotFindFavouriteForAsset
		}

		result = append(
			result,
			AudienceFavourite{
				Id:          favourite.Id,
				Description: favourite.Description,
				Info:        audience,
			},
		)
	}

	return &result, nil
}
