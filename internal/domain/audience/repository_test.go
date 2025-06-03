package audience_test

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/domain/audience"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByIds(t *testing.T) {
	t.Run("should return all audiences when all IDs exist", func(t *testing.T) {
		// Arrange
		aud1ID := uuid.New()
		aud2ID := uuid.New()

		mockDB := &database.IMDatabase{
			AudienceStorage: map[uuid.UUID]database.IMAudienceModel{
				aud1ID: {
					Id:                 aud1ID,
					Gender:             "Female",
					BirthCountry:       "Canada",
					AgeGroup:           "18-24",
					SocialMediaHours:   5.5,
					PurchasesLastMonth: 2,
				},
				aud2ID: {
					Id:                 aud2ID,
					Gender:             "Male",
					BirthCountry:       "USA",
					AgeGroup:           "25-34",
					SocialMediaHours:   3.0,
					PurchasesLastMonth: 4,
				},
			},
		}

		repo := audience.NewInMemoryDBAudienceRepository(mockDB)

		expectedResult := []audience.Audience{
			{
				Id:                 aud1ID,
				Gender:             "Female",
				BirthCountry:       "Canada",
				AgeGroup:           "18-24",
				SocialMediaHours:   5.5,
				PurchasesLastMonth: 2,
			},
			{
				Id:                 aud2ID,
				Gender:             "Male",
				BirthCountry:       "USA",
				AgeGroup:           "25-34",
				SocialMediaHours:   3.0,
				PurchasesLastMonth: 4,
			},
		}

		// Act
		actual, _ := repo.GetByIds([]uuid.UUID{aud1ID, aud2ID})

		// Assert
		assert.Equal(t, expectedResult, actual)
	})

	t.Run("should return only matching audiences when some IDs do not exist", func(t *testing.T) {
		// Arrange
		aud1ID := uuid.New()
		missingID := uuid.New()

		mockDB := &database.IMDatabase{
			AudienceStorage: map[uuid.UUID]database.IMAudienceModel{
				aud1ID: {
					Id:                 aud1ID,
					Gender:             "Female",
					BirthCountry:       "Brazil",
					AgeGroup:           "35-44",
					SocialMediaHours:   2.0,
					PurchasesLastMonth: 1,
				},
			},
		}

		repo := audience.NewInMemoryDBAudienceRepository(mockDB)

		expectedResult := []audience.Audience{
			{
				Id:                 aud1ID,
				Gender:             "Female",
				BirthCountry:       "Brazil",
				AgeGroup:           "35-44",
				SocialMediaHours:   2.0,
				PurchasesLastMonth: 1,
			},
		}

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{aud1ID, missingID})

		// Assert
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should return empty slice when input is empty", func(t *testing.T) {
		// Arrange
		mockDB := &database.IMDatabase{
			AudienceStorage: map[uuid.UUID]database.IMAudienceModel{},
		}
		repo := audience.NewInMemoryDBAudienceRepository(mockDB)

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{})

		// Assert
		assert.Empty(t, result)
	})
}

func TestGetById(t *testing.T) {
	t.Run("should return audience when ID exists", func(t *testing.T) {
		// Arrange
		audID := uuid.New()
		mockDB := &database.IMDatabase{
			AudienceStorage: map[uuid.UUID]database.IMAudienceModel{
				audID: {
					Id:                 audID,
					Gender:             "Other",
					BirthCountry:       "Germany",
					AgeGroup:           "45-54",
					SocialMediaHours:   1.2,
					PurchasesLastMonth: 0,
				},
			},
		}

		repo := audience.NewInMemoryDBAudienceRepository(mockDB)

		expected := &audience.Audience{
			Id:                 audID,
			Gender:             "Other",
			BirthCountry:       "Germany",
			AgeGroup:           "45-54",
			SocialMediaHours:   1.2,
			PurchasesLastMonth: 0,
		}

		// Act
		result, err := repo.GetById(audID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return error when ID does not exist", func(t *testing.T) {
		// Arrange
		mockDB := &database.IMDatabase{
			AudienceStorage: map[uuid.UUID]database.IMAudienceModel{},
		}
		repo := audience.NewInMemoryDBAudienceRepository(mockDB)

		missingID := uuid.New()

		// Act
		result, err := repo.GetById(missingID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
