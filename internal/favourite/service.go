package favourite

import (
	"fmt"
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/insight"
	"platform-go-challenge/internal/utils"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type FavouriteService interface {
	GetPaginatedForUser(UserId uuid.UUID, pageSize int, pageNumber int) (*AssetFavourites, *utils.Pagination, error)
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

	fmt.Println(chartIds)

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
