package favourite_test

import (
	"errors"
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/favourite"
	"platform-go-challenge/internal/insight"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockFavouriteRepo struct {
	getByUserIdPaginatedFn func(userId uuid.UUID, pageSize, pageNumber int) ([]favourite.Favourite, utils.Pagination, error)
	createFn               func(fav favourite.Favourite) (*favourite.Favourite, error)
}

func (m *mockFavouriteRepo) GetByUserIdPaginated(userId uuid.UUID, pageSize, pageNumber int) ([]favourite.Favourite, utils.Pagination, error) {
	return m.getByUserIdPaginatedFn(userId, pageSize, pageNumber)
}

func (m *mockFavouriteRepo) Create(fav favourite.Favourite) (*favourite.Favourite, error) {
	return m.createFn(fav)
}

type mockChartRepo struct {
	getByIdsFn func(ids uuid.UUIDs) ([]chart.Chart, error)
	getByIdFn  func(id uuid.UUID) (*chart.Chart, error)
}

func (m *mockChartRepo) GetByIds(ids uuid.UUIDs) ([]chart.Chart, error) {
	return m.getByIdsFn(ids)
}

func (m *mockChartRepo) GetById(id uuid.UUID) (*chart.Chart, error) {
	return m.getByIdFn(id)
}

type mockInsightRepo struct {
	getByIdsFn func(ids uuid.UUIDs) ([]insight.Insight, error)
	getByIdFn  func(id uuid.UUID) (*insight.Insight, error)
}

func (m *mockInsightRepo) GetByIds(ids uuid.UUIDs) ([]insight.Insight, error) {
	return m.getByIdsFn(ids)
}

func (m *mockInsightRepo) GetById(id uuid.UUID) (*insight.Insight, error) {
	return m.getByIdFn(id)
}

type mockAudienceRepo struct {
	getByIdsFn func(ids uuid.UUIDs) ([]audience.Audience, error)
	getByIdFn  func(id uuid.UUID) (*audience.Audience, error)
}

func (m *mockAudienceRepo) GetByIds(ids uuid.UUIDs) ([]audience.Audience, error) {
	return m.getByIdsFn(ids)
}

func (m *mockAudienceRepo) GetById(id uuid.UUID) (*audience.Audience, error) {
	return m.getByIdFn(id)
}

func TestShouldReturnPaginatedFavouritesWhenGetPaginatedForUser(t *testing.T) {
	// Arrange
	userId := uuid.New()
	pageSize := 10
	pageNumber := 1

	favourites := []favourite.Favourite{
		{Id: uuid.New(), UserId: userId, AssetId: uuid.New(), AssetType: favourite.AssetTypeChart},
		{Id: uuid.New(), UserId: userId, AssetId: uuid.New(), AssetType: favourite.AssetTypeInsight},
		{Id: uuid.New(), UserId: userId, AssetId: uuid.New(), AssetType: favourite.AssetTypeAudience},
	}
	pagination := utils.Pagination{Page: pageNumber, PageSize: pageSize, MaxPage: 3}

	mockFavRepo := &mockFavouriteRepo{
		getByUserIdPaginatedFn: func(uId uuid.UUID, ps, pn int) ([]favourite.Favourite, utils.Pagination, error) {
			assert.Equal(t, userId, uId)
			assert.Equal(t, pageSize, ps)
			assert.Equal(t, pageNumber, pn)
			return favourites, pagination, nil
		},
	}

	mockChartRepo := &mockChartRepo{
		getByIdsFn: func(ids uuid.UUIDs) ([]chart.Chart, error) {
			assert.Len(t, ids, 1)
			return []chart.Chart{{Id: ids[0]}}, nil
		},
	}

	mockInsightRepo := &mockInsightRepo{
		getByIdsFn: func(ids uuid.UUIDs) ([]insight.Insight, error) {
			assert.Len(t, ids, 1)
			return []insight.Insight{{Id: ids[0]}}, nil
		},
	}

	mockAudienceRepo := &mockAudienceRepo{
		getByIdsFn: func(ids uuid.UUIDs) ([]audience.Audience, error) {
			assert.Len(t, ids, 1)
			return []audience.Audience{{Id: ids[0]}}, nil
		},
	}

	service := favourite.NewFavouriteService(favourite.FavouriteServiceDependencies{
		FavouriteRepository: mockFavRepo,
		ChartRepository:     mockChartRepo,
		InsightRepository:   mockInsightRepo,
		AudienceRepository:  mockAudienceRepo,
	})

	// Act
	result, pag, err := service.GetPaginatedForUser(userId, pageSize, pageNumber)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &pagination, pag)
	assert.Len(t, result.Charts, 1)
	assert.Len(t, result.Insights, 1)
	assert.Len(t, result.Audiences, 1)
}

func TestShouldCreateFavouriteWhenAssetExists(t *testing.T) {
	// Arrange
	userId := uuid.New()
	assetId := uuid.New()
	description := "desc"

	mockChartRepo := &mockChartRepo{
		getByIdFn: func(id uuid.UUID) (*chart.Chart, error) {
			if id == assetId {
				return &chart.Chart{Id: id}, nil
			}
			return nil, nil
		},
	}
	mockInsightRepo := &mockInsightRepo{
		getByIdFn: func(id uuid.UUID) (*insight.Insight, error) { return nil, nil },
	}
	mockAudienceRepo := &mockAudienceRepo{
		getByIdFn: func(id uuid.UUID) (*audience.Audience, error) { return nil, nil },
	}
	mockFavRepo := &mockFavouriteRepo{
		createFn: func(fav favourite.Favourite) (*favourite.Favourite, error) {
			assert.Equal(t, userId, fav.UserId)
			assert.Equal(t, assetId, fav.AssetId)
			assert.Equal(t, favourite.AssetTypeChart, fav.AssetType)
			assert.Equal(t, description, fav.Description)
			created := fav
			created.Id = uuid.New()
			return &created, nil
		},
	}

	service := favourite.NewFavouriteService(favourite.FavouriteServiceDependencies{
		FavouriteRepository: mockFavRepo,
		ChartRepository:     mockChartRepo,
		InsightRepository:   mockInsightRepo,
		AudienceRepository:  mockAudienceRepo,
	})

	// Act
	created, err := service.CreateForUser(userId, assetId, description)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, userId, created.UserId)
	assert.Equal(t, assetId, created.AssetId)
	assert.Equal(t, description, created.Description)
	assert.Equal(t, favourite.AssetTypeChart, created.AssetType)
}

func TestShouldReturnAssetNotFoundWhenCreateForUserAndAssetDoesNotExist(t *testing.T) {
	// Arrange
	userId := uuid.New()
	assetId := uuid.New()
	mockChartRepo := &mockChartRepo{
		getByIdFn: func(id uuid.UUID) (*chart.Chart, error) { return nil, nil },
	}
	mockInsightRepo := &mockInsightRepo{
		getByIdFn: func(id uuid.UUID) (*insight.Insight, error) { return nil, nil },
	}
	mockAudienceRepo := &mockAudienceRepo{
		getByIdFn: func(id uuid.UUID) (*audience.Audience, error) { return nil, nil },
	}
	mockFavRepo := &mockFavouriteRepo{}

	service := favourite.NewFavouriteService(favourite.FavouriteServiceDependencies{
		FavouriteRepository: mockFavRepo,
		ChartRepository:     mockChartRepo,
		InsightRepository:   mockInsightRepo,
		AudienceRepository:  mockAudienceRepo,
	})

	// Act
	created, err := service.CreateForUser(userId, assetId, "desc")

	// Assert
	assert.Nil(t, created)
	assert.ErrorIs(t, err, favourite.ErrAssetNotFound)
}

func TestShouldReturnErrorWhenCreateForUserFailsToSave(t *testing.T) {
	// Arrange
	userId := uuid.New()
	assetId := uuid.New()
	mockChartRepo := &mockChartRepo{
		getByIdFn: func(id uuid.UUID) (*chart.Chart, error) {
			if id == assetId {
				return &chart.Chart{Id: id}, nil
			}
			return nil, nil
		},
	}
	mockInsightRepo := &mockInsightRepo{
		getByIdFn: func(id uuid.UUID) (*insight.Insight, error) { return nil, nil },
	}
	mockAudienceRepo := &mockAudienceRepo{
		getByIdFn: func(id uuid.UUID) (*audience.Audience, error) { return nil, nil },
	}
	mockFavRepo := &mockFavouriteRepo{
		createFn: func(fav favourite.Favourite) (*favourite.Favourite, error) {
			return nil, errors.New("db error")
		},
	}

	service := favourite.NewFavouriteService(favourite.FavouriteServiceDependencies{
		FavouriteRepository: mockFavRepo,
		ChartRepository:     mockChartRepo,
		InsightRepository:   mockInsightRepo,
		AudienceRepository:  mockAudienceRepo,
	})

	// Act
	created, err := service.CreateForUser(userId, assetId, "desc")

	// Assert
	assert.Nil(t, created)
	assert.ErrorIs(t, err, favourite.ErrCouldNotSaveFavourite)
}
