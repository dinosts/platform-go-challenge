package favourite

import (
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/insight"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type FavouriteService interface {
	GetPaginatedForUser(UserId uuid.UUID)
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

func NewUserService(dependencies FavouriteServiceDependencies) favouriteService {
	return favouriteService{
		Dependencies: dependencies,
	}
}

func (service *favouriteService) GetPaginatedForUser(UserId uuid.UUID, pageSize int, pageNumber int) {
	favourites, _ := service.Dependencies.FavouriteRepository.GetByUserIdPaginated(UserId, pageSize, pageNumber)

	groupedFavourites := groupFavouritesByAssetType(favourites)

	chartIds := extractIds(groupedFavourites.Charts)
	insightIds := extractIds(groupedFavourites.Insights)
	audienceIds := extractIds(groupedFavourites.Audience)

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
		return nil, err
	}

	print(charts, insights, audiences)
}
