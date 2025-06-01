package favourite_test

import (
	"errors"
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/favourite"
	"platform-go-challenge/internal/insight"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildAssetFavourites(t *testing.T) {
	t.Run("should build grouped asset favourites successfully", func(t *testing.T) {
		// Arrange
		chartID := uuid.New()
		insightID := uuid.New()
		audienceID := uuid.New()
		favID1 := uuid.New()
		favID2 := uuid.New()
		favID3 := uuid.New()

		charts := []chart.Chart{{Id: chartID}}
		insights := []insight.Insight{{Id: insightID}}
		audiences := []audience.Audience{{Id: audienceID}}

		favs := []favourite.Favourite{
			{Id: favID1, AssetId: chartID, AssetType: "charts", Description: "Chart A"},
			{Id: favID2, AssetId: insightID, AssetType: "insights", Description: "Insight B"},
			{Id: favID3, AssetId: audienceID, AssetType: "audiences", Description: "Audience C"},
		}

		// Act
		result, err := favourite.BuildAssetFavourites(favs, charts, insights, audiences)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)

		assert.Len(t, result.Charts, 1)
		assert.Equal(t, favID1, result.Charts[0].Id)
		assert.Equal(t, chartID, result.Charts[0].Info.Id)

		assert.Len(t, result.Insights, 1)
		assert.Equal(t, favID2, result.Insights[0].Id)
		assert.Equal(t, insightID, result.Insights[0].Info.Id)

		assert.Len(t, result.Audiences, 1)
		assert.Equal(t, favID3, result.Audiences[0].Id)
		assert.Equal(t, audienceID, result.Audiences[0].Info.Id)
	})

	t.Run("should fail if favourite is missing for a provided asset", func(t *testing.T) {
		// Arrange
		chartID := uuid.New()
		charts := []chart.Chart{{Id: chartID}}
		insights := []insight.Insight{}
		audiences := []audience.Audience{}
		favs := []favourite.Favourite{} // no favourites provided

		// Act
		_, err := favourite.BuildAssetFavourites(favs, charts, insights, audiences)

		// Assert
		assert.Error(t, err)
		assert.True(t, errors.Is(err, favourite.ErrCouldNotFindFavouriteForAsset))
	})

	t.Run("should return empty grouped favourites when asset slices are empty", func(t *testing.T) {
		// Arrange
		favs := []favourite.Favourite{}
		charts := []chart.Chart{}
		insights := []insight.Insight{}
		audiences := []audience.Audience{}

		// Act
		result, err := favourite.BuildAssetFavourites(favs, charts, insights, audiences)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Charts, 0)
		assert.Len(t, result.Insights, 0)
		assert.Len(t, result.Audiences, 0)
	})
}
