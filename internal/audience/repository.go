package audience

import (
	"platform-go-challenge/internal/database"

	"github.com/google/uuid"
)

type AudienceRepository interface {
	GetByIds(ids uuid.UUIDs) ([]Audience, error)
	GetById(id uuid.UUID) (*Audience, error)
}

type inMemoryDBAudienceRepository struct {
	DB *database.InMemoryDatabase
}

func NewInMemoryDBAudienceRepository(db *database.InMemoryDatabase) *inMemoryDBAudienceRepository {
	return &inMemoryDBAudienceRepository{
		DB: db,
	}
}

func InMemoryDBAudienceModelToDTO(model database.AudienceModel) Audience {
	return Audience{
		Id:                 model.Id,
		Gender:             model.Gender,
		BirthCountry:       model.BirthCountry,
		AgeGroup:           model.AgeGroup,
		SocialMediaHours:   model.SocialMediaHours,
		PurchasesLastMonth: model.PurchasesLastMonth,
	}
}

func (repo *inMemoryDBAudienceRepository) GetByIds(ids uuid.UUIDs) ([]Audience, error) {
	result := []Audience{}
	for _, id := range ids {
		model, found := repo.DB.AudienceStorage[id]
		if !found {
			continue
		}
		dto := InMemoryDBAudienceModelToDTO(model)
		result = append(result, dto)
	}
	return result, nil
}

func (repo *inMemoryDBAudienceRepository) GetById(id uuid.UUID) (*Audience, error) {
	audience, err := database.IMStorageGetById(
		id,
		repo.DB.AudienceStorage,
	)
	if err != nil {
		return nil, err
	}

	dto := InMemoryDBAudienceModelToDTO(*audience)

	return &dto, nil
}
