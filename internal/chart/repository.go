package chart

import (
	"platform-go-challenge/internal/database"

	"github.com/google/uuid"
)

type ChartRepository interface {
	GetByIds(ids uuid.UUIDs) ([]Chart, error)
	GetById(id uuid.UUID) (*Chart, error)
}

type inMemoryDBChartRepository struct {
	DB *database.InMemoryDatabase
}

func NewInMemoryDBChartRepository(db *database.InMemoryDatabase) *inMemoryDBChartRepository {
	return &inMemoryDBChartRepository{
		DB: db,
	}
}

func InMemoryDBChartModelToDTO(chartModel database.ChartModel) Chart {
	return Chart{
		Id:         chartModel.Id,
		Title:      chartModel.Title,
		XAxisTitle: chartModel.XAxisTitle,
		YAxisTitle: chartModel.YAxisTitle,
		Data:       chartModel.Data,
	}
}

func (repo *inMemoryDBChartRepository) GetByIds(ids uuid.UUIDs) ([]Chart, error) {
	result := []Chart{}

	for _, id := range ids {
		v, found := repo.DB.ChartStorage[id]
		if !found {
			continue
		}

		dto := InMemoryDBChartModelToDTO(v)

		result = append(result, dto)
	}

	return result, nil
}

func (repo *inMemoryDBChartRepository) GetById(id uuid.UUID) (*Chart, error) {
	chart, err := database.IMStorageGetById(
		id,
		repo.DB.ChartStorage,
	)
	if err != nil {
		return nil, err
	}

	dto := InMemoryDBChartModelToDTO(*chart)

	return &dto, nil
}
