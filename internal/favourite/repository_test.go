package favourite_test

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/favourite"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByUserIdPaginated(t *testing.T) {
	// Setup some users and favourites
	user1 := uuid.New()
	user2 := uuid.New()

	fav1 := database.FavouriteModel{
		Id:          uuid.New(),
		UserId:      user1,
		AssetId:     uuid.New(),
		AssetType:   "charts",
		Description: "fav1",
	}
	fav2 := database.FavouriteModel{
		Id:          uuid.New(),
		UserId:      user1,
		AssetId:     uuid.New(),
		AssetType:   "charts",
		Description: "fav2",
	}
	fav3 := database.FavouriteModel{
		Id:          uuid.New(),
		UserId:      user2,
		AssetId:     uuid.New(),
		AssetType:   "insights",
		Description: "fav3",
	}
	fav4 := database.FavouriteModel{
		Id:          uuid.New(),
		UserId:      user1,
		AssetId:     uuid.New(),
		AssetType:   "audiences",
		Description: "fav4",
	}

	t.Run("returns paginated favourites for user", func(t *testing.T) {
		// Arrange
		storage := make(map[uuid.UUID]database.FavouriteModel)
		storage[fav1.Id] = fav1
		storage[fav2.Id] = fav2
		storage[fav3.Id] = fav3
		storage[fav4.Id] = fav4

		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{
			FavouriteStorage: storage,
		})

		pageSize := 2
		pageNumber := 0

		expectedItemsLen := pageSize
		expectedMaxPage := 1

		// Act
		result, _ := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.Equal(t, pageNumber, result.Page)
		assert.Equal(t, pageSize, result.PageSize)

		assert.Equal(t, expectedItemsLen, len(result.Items))
		assert.Equal(t, expectedMaxPage, result.MaxPage)

		for _, fav := range result.Items {
			assert.Equal(t, user1, fav.UserId)
		}
	})

	t.Run("returns second page correctly", func(t *testing.T) {
		storage := make(map[uuid.UUID]database.FavouriteModel)
		storage[fav1.Id] = fav1
		storage[fav2.Id] = fav2
		storage[fav3.Id] = fav3
		storage[fav4.Id] = fav4

		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{
			FavouriteStorage: storage,
		})

		pageSize := 2
		pageNumber := 1

		expectedItemsLen := 1
		expectedMaxPage := 1

		// Arrange
		result, _ := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.Equal(t, expectedItemsLen, len(result.Items))
		assert.Equal(t, expectedMaxPage, result.MaxPage)

		assert.Equal(t, pageNumber, result.Page)
		assert.Equal(t, pageSize, result.PageSize)

		for _, fav := range result.Items {
			assert.Equal(t, user1, fav.UserId)
		}
	})

	t.Run("returns empty when page out of range", func(t *testing.T) {
		// Arrange
		storage := make(map[uuid.UUID]database.FavouriteModel)
		storage[fav1.Id] = fav1

		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{
			FavouriteStorage: storage,
		})

		pageSize := 1
		pageNumber := 5

		// Act
		result, _ := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.Empty(t, result.Items)

		assert.Equal(t, pageNumber, result.Page)
		assert.Equal(t, pageSize, result.PageSize)

		assert.Equal(t, 1, result.MaxPage)
	})

	t.Run("returns empty on no favourites", func(t *testing.T) {
		// Arrange
		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{
			FavouriteStorage: map[uuid.UUID]database.FavouriteModel{},
		})

		// Act
		result, _ := repo.GetByUserIdPaginated(user1, 10, 0)

		// Assert
		assert.Empty(t, result.Items)
		assert.Equal(t, 0, result.MaxPage)
	})
}
