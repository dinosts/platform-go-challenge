package favourite_test

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/favourite"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByUserIdPaginated(t *testing.T) {
	user1 := uuid.New()
	user2 := uuid.New()

	fav1 := database.FavouriteModel{
		Id: uuid.New(), UserId: user1, AssetId: uuid.New(), AssetType: "charts", Description: "fav1",
	}
	fav2 := database.FavouriteModel{
		Id: uuid.New(), UserId: user1, AssetId: uuid.New(), AssetType: "charts", Description: "fav2",
	}
	fav3 := database.FavouriteModel{
		Id: uuid.New(), UserId: user2, AssetId: uuid.New(), AssetType: "insights", Description: "fav3",
	}
	fav4 := database.FavouriteModel{
		Id: uuid.New(), UserId: user1, AssetId: uuid.New(), AssetType: "audiences", Description: "fav4",
	}

	t.Run("should return first page of favourites for a user", func(t *testing.T) {
		// Arrange
		storage := map[uuid.UUID]database.FavouriteModel{
			fav1.Id: fav1, fav2.Id: fav2, fav3.Id: fav3, fav4.Id: fav4,
		}
		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{FavouriteStorage: storage})
		pageSize := 2
		pageNumber := 0
		expectedItemsLen := 2
		expectedMaxPage := 1

		// Act
		result, pagination, err := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedItemsLen, len(result))
		assert.Equal(t, expectedMaxPage, pagination.MaxPage)
		assert.Equal(t, pageSize, pagination.PageSize)
		assert.Equal(t, pageNumber, pagination.Page)
		for _, fav := range result {
			assert.Equal(t, user1, fav.UserId)
		}
	})

	t.Run("should return second page of favourites for a user", func(t *testing.T) {
		// Arrange
		storage := map[uuid.UUID]database.FavouriteModel{
			fav1.Id: fav1, fav2.Id: fav2, fav3.Id: fav3, fav4.Id: fav4,
		}
		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{FavouriteStorage: storage})
		pageSize := 2
		pageNumber := 1
		expectedItemsLen := 1
		expectedMaxPage := 1

		// Act
		result, pagination, err := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedItemsLen, len(result))
		assert.Equal(t, expectedMaxPage, pagination.MaxPage)
		assert.Equal(t, pageSize, pagination.PageSize)
		assert.Equal(t, pageNumber, pagination.Page)
		for _, fav := range result {
			assert.Equal(t, user1, fav.UserId)
		}
	})

	t.Run("should return empty slice when page number is out of range", func(t *testing.T) {
		// Arrange
		storage := map[uuid.UUID]database.FavouriteModel{fav1.Id: fav1}
		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{FavouriteStorage: storage})
		pageSize := 1
		pageNumber := 5

		// Act
		result, pagination, err := repo.GetByUserIdPaginated(user1, pageSize, pageNumber)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, 0, pagination.MaxPage)
		assert.Equal(t, pageSize, pagination.PageSize)
		assert.Equal(t, pageNumber, pagination.Page)
	})

	t.Run("should return empty slice when user has no favourites", func(t *testing.T) {
		// Arrange
		repo := favourite.NewInMemoryDBFavouriteRepository(&database.InMemoryDatabase{
			FavouriteStorage: map[uuid.UUID]database.FavouriteModel{},
		})

		// Act
		result, pagination, err := repo.GetByUserIdPaginated(user1, 10, 0)

		// Assert
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, 0, pagination.MaxPage)
	})
}

func TestCreate(t *testing.T) {
	t.Run("should store favourite in memory database", func(t *testing.T) {
		// Arrange
		db := &database.InMemoryDatabase{FavouriteStorage: make(map[uuid.UUID]database.FavouriteModel)}
		repo := favourite.NewInMemoryDBFavouriteRepository(db)

		newFav := favourite.Favourite{
			Id: uuid.New(), UserId: uuid.New(), AssetId: uuid.New(),
			AssetType: "charts", Description: "created fav",
		}

		// Act
		created, err := repo.Create(newFav)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, &newFav, created)

		stored, ok := db.FavouriteStorage[newFav.Id]
		assert.True(t, ok)
		assert.Equal(t, newFav.Id, stored.Id)
		assert.Equal(t, newFav.UserId, stored.UserId)
		assert.Equal(t, newFav.AssetId, stored.AssetId)
		assert.Equal(t, string(newFav.AssetType), stored.AssetType)
		assert.Equal(t, newFav.Description, stored.Description)
	})
}
