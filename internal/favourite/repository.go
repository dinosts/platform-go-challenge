package favourite

import (
	"maps"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/utils"
	"slices"
	"sort"

	"github.com/google/uuid"
)

type FavouriteRepository interface {
	GetByUserIdPaginated(userId uuid.UUID, pageSize int, pageNumber int) ([]Favourite, utils.Pagination, error)
	GetById(id uuid.UUID) (*Favourite, error)
	Create(favourite Favourite) (*Favourite, error)
	Update(favourite Favourite) (*Favourite, error)
	Delete(id uuid.UUID) error
}

type inMemoryDBFavouriteRepository struct {
	DB *database.InMemoryDatabase
}

func NewInMemoryDBFavouriteRepository(db *database.InMemoryDatabase) *inMemoryDBFavouriteRepository {
	return &inMemoryDBFavouriteRepository{
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

func DTOToInMemoryDBFavouriteModel(dto Favourite) database.FavouriteModel {
	return database.FavouriteModel{
		Id:          dto.Id,
		UserId:      dto.UserId,
		AssetId:     dto.AssetId,
		AssetType:   string(dto.AssetType),
		Description: dto.Description,
	}
}

func (repo *inMemoryDBFavouriteRepository) GetById(id uuid.UUID) (*Favourite, error) {
	favourite, err := database.IMStorageGetById(id, repo.DB.FavouriteStorage)
	if err != nil {
		return nil, err
	}

	dto := InMemoryDBFavouriteModelToDTO(*favourite)

	return &dto, nil
}

func (repo *inMemoryDBFavouriteRepository) GetByUserIdPaginated(userId uuid.UUID, pageSize int, pageNumber int) ([]Favourite, utils.Pagination, error) {
	var result []Favourite

	totalCount := 0
	offset := pageSize * pageNumber

	sortedFavs := slices.Collect(maps.Values(repo.DB.FavouriteStorage))

	sort.Slice(
		sortedFavs,
		func(i, j int) bool { return sortedFavs[i].Id.String() < sortedFavs[j].Id.String() },
	)

	for _, fav := range sortedFavs {
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

	maxPage := utils.CalculateMaxPages(totalCount, pageSize)

	return result, utils.Pagination{Page: pageNumber, PageSize: pageSize, MaxPage: maxPage}, nil
}

func (repo *inMemoryDBFavouriteRepository) Create(favourite Favourite) (*Favourite, error) {
	repo.DB.FavouriteStorage[favourite.Id] = DTOToInMemoryDBFavouriteModel(favourite)
	return &favourite, nil
}

func (repo *inMemoryDBFavouriteRepository) Update(favourite Favourite) (*Favourite, error) {
	repo.DB.FavouriteStorage[favourite.Id] = DTOToInMemoryDBFavouriteModel(favourite)
	return &favourite, nil
}

func (repo *inMemoryDBFavouriteRepository) Delete(id uuid.UUID) error {
	favourite, err := database.IMStorageGetById(id, repo.DB.FavouriteStorage)
	if err != nil {
		return ErrFavouriteNotFound
	}

	delete(repo.DB.FavouriteStorage, favourite.Id)

	return nil
}
