package insight

import (
	"platform-go-challenge/internal/database"

	"github.com/google/uuid"
)

type InsightRepository interface {
	GetByIds(ids uuid.UUIDs) ([]Insight, error)
	GetById(id uuid.UUID) (*Insight, error)
}

type inMemoryDBInsightRepository struct {
	DB *database.InMemoryDatabase
}

func NewInMemoryDBInsightRepository(db *database.InMemoryDatabase) *inMemoryDBInsightRepository {
	return &inMemoryDBInsightRepository{
		DB: db,
	}
}

func InMemoryDBInsightModelToDTO(insightModel database.InsightModel) Insight {
	return Insight{
		Id:   insightModel.Id,
		Text: insightModel.Text,
	}
}

func (repo *inMemoryDBInsightRepository) GetByIds(ids uuid.UUIDs) ([]Insight, error) {
	result := []Insight{}
	for _, id := range ids {
		v, found := repo.DB.InsightStorage[id]

		if !found {
			continue
		}

		dto := InMemoryDBInsightModelToDTO(v)

		result = append(result, dto)
	}

	return result, nil
}

func (repo *inMemoryDBInsightRepository) GetById(id uuid.UUID) (*Insight, error) {
	insight, err := database.IMStorageGetById(
		id,
		repo.DB.InsightStorage,
	)
	if err != nil {
		return nil, err
	}

	dto := InMemoryDBInsightModelToDTO(*insight)

	return &dto, nil
}
