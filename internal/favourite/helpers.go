package favourite

import "github.com/google/uuid"

func groupFavouritesByAssetType(favourites []Favourite) GroupedFavourites {
	grouped := GroupedFavourites{}

	for _, fav := range favourites {
		switch fav.AssetType {
		case AssetTypeChart:
			grouped.Charts = append(grouped.Charts, fav)
		case AssetTypeInsight:
			grouped.Insights = append(grouped.Insights, fav)
		case AssetTypeAudience:
			grouped.Audience = append(grouped.Audience, fav)
		}
	}

	return grouped
}

func extractIds(favourites []Favourite) uuid.UUIDs {
	result := make([]uuid.UUID, len(favourites))

	for i, fav := range favourites {
		result[i] = fav.Id
	}

	return result
}
