package favourite

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/utils"

	"github.com/google/uuid"
)

type FavouriteRepository interface {
	GetByUserIdPaginated(userId uuid.UUID, pageSize int, pageNumber int) ([]Favourite, utils.Pagination, error)
}

type inMemoryDBFavouriteRepository struct {
	DB *database.InMemoryDatabase
}

func NewInMemoryDBFavouriteRepository(db *database.InMemoryDatabase) inMemoryDBFavouriteRepository {
	return inMemoryDBFavouriteRepository{
		DB: db,
	}
}

func InMemoryDBFavouriteModelToDTO(model database.FavouriteModel) Favourite {
	return Favourite{
		Id:          model.Id,
		UserId:      model.UserId,
		AssetId:     model.AssetId,
		AssetType:   AssetType(model.AssetType),
		Description: model.Description,
	}
}

func (repo *inMemoryDBFavouriteRepository) GetByUserIdPaginated(userId uuid.UUID, pageSize int, pageNumber int) ([]Favourite, utils.Pagination, error) {
	var result []Favourite

	totalCount := 0
	offset := pageSize * pageNumber

	for _, fav := range repo.DB.FavouriteStorage {
		if fav.UserId != userId {
			continue
		}

		totalCount++

		if totalCount <= offset {
			continue
		}

		if len(result) != pageSize {
			result = append(result, InMemoryDBFavouriteModelToDTO(fav))
		}
	}

	// int casting rounds down. (wanted behaviour cause we start pages from 0)
	totalPages := int(totalCount / pageSize)

	return result, utils.Pagination{Page: pageNumber, PageSize: pageSize, MaxPage: totalPages}, nil
}
