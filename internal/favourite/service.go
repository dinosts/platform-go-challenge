package favourite

import (
	"errors"
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/insight"
	"platform-go-challenge/internal/utils"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type FavouriteService interface {
	GetPaginatedForUser(UserId uuid.UUID, pageSize int, pageNumber int) (*AssetFavourites, *utils.Pagination, error)
	CreateForUser(UserId uuid.UUID, assetId uuid.UUID, description string) (*Favourite, error)
	Update(userId uuid.UUID, favouriteId uuid.UUID, newDescription string) (*Favourite, error)
}

type FavouriteServiceDependencies struct {
	FavouriteRepository FavouriteRepository
	ChartRepository     chart.ChartRepository
	InsightRepository   insight.InsightRepository
	AudienceRepository  audience.AudienceRepository
}

type favouriteService struct {
	Dependencies FavouriteServiceDependencies
}

func NewFavouriteService(dependencies FavouriteServiceDependencies) favouriteService {
	return favouriteService{
		Dependencies: dependencies,
	}
}

func (service *favouriteService) GetPaginatedForUser(UserId uuid.UUID, pageSize int, pageNumber int) (*AssetFavourites, *utils.Pagination, error) {
	favourites, pagination, err := service.Dependencies.FavouriteRepository.GetByUserIdPaginated(UserId, pageSize, pageNumber)
	if err != nil {
		return nil, nil, err
	}

	chartIds := ExtractAssetTypeIds(AssetTypeChart, favourites)
	insightIds := ExtractAssetTypeIds(AssetTypeInsight, favourites)
	audienceIds := ExtractAssetTypeIds(AssetTypeAudience, favourites)

	var (
		charts    []chart.Chart
		insights  []insight.Insight
		audiences []audience.Audience
	)

	g := new(errgroup.Group)

	g.Go(func() error {
		var err error
		charts, err = service.Dependencies.ChartRepository.GetByIds(chartIds)
		return err
	})

	g.Go(func() error {
		var err error
		insights, err = service.Dependencies.InsightRepository.GetByIds(insightIds)
		return err
	})

	g.Go(func() error {
		var err error
		audiences, err = service.Dependencies.AudienceRepository.GetByIds(audienceIds)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	result, err := BuildAssetFavourites(favourites, charts, insights, audiences)
	if err != nil {
		return nil, nil, err
	}

	return result, &pagination, nil
}

func (service *favouriteService) detectAssetType(assetId uuid.UUID) (AssetType, error) {
	var (
		chart    *chart.Chart
		insight  *insight.Insight
		audience *audience.Audience
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		chart, _ = service.Dependencies.ChartRepository.GetById(assetId)
	}()
	go func() {
		defer wg.Done()
		insight, _ = service.Dependencies.InsightRepository.GetById(assetId)
	}()
	go func() {
		defer wg.Done()
		audience, _ = service.Dependencies.AudienceRepository.GetById(assetId)
	}()

	wg.Wait()

	switch {
	case chart != nil:
		return AssetTypeChart, nil
	case insight != nil:
		return AssetTypeInsight, nil
	case audience != nil:
		return AssetTypeAudience, nil
	default:
		return "", ErrAssetNotFound
	}
}

func (service *favouriteService) CreateForUser(userId uuid.UUID, assetId uuid.UUID, description string) (*Favourite, error) {
	assetType, err := service.detectAssetType(assetId)
	if err != nil {
		return nil, err
	}

	favourite := Favourite{
		Id:          uuid.New(),
		UserId:      userId,
		AssetId:     assetId,
		AssetType:   assetType,
		Description: description,
	}

	fav, err := service.Dependencies.FavouriteRepository.Create(favourite)
	if err != nil {
		return nil, ErrCouldNotSaveFavourite
	}

	return fav, nil
}

func (service *favouriteService) Update(userId uuid.UUID, favouriteId uuid.UUID, newDescription string) (*Favourite, error) {
	favourite, err := service.Dependencies.FavouriteRepository.GetById(favouriteId)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return nil, ErrFavouriteNotFound
		}

		return nil, utils.ErrUnexpected
	}

	if userId != favourite.UserId {
		return nil, ErrFavouriteNotUnderGivenUser
	}

	if newDescription != "" {
		favourite.Description = newDescription
	}

	return service.Dependencies.FavouriteRepository.Update(*favourite)
}
