package insight_test

import (
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/domain/insight"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDBInsightRepository_GetByIds(t *testing.T) {
	t.Run("should return all matching insights when valid ids are provided", func(t *testing.T) {
		// Arrange
		id1 := uuid.New()
		id2 := uuid.New()

		mockInsight1 := database.IMInsightModel{Id: id1, Text: "Insight 1"}
		mockInsight2 := database.IMInsightModel{Id: id2, Text: "Insight 2"}

		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{
				id1: mockInsight1,
				id2: mockInsight2,
			},
		}

		repo := insight.NewInMemoryDBInsightRepository(db)

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{id1, id2})

		// Assert
		assert.Len(t, result, 2)
		assert.Contains(t, result, insight.Insight{Id: id1, Text: "Insight 1"})
		assert.Contains(t, result, insight.Insight{Id: id2, Text: "Insight 2"})
	})

	t.Run("should return only matching insights and ignore missing ones", func(t *testing.T) {
		// Arrange
		idExisting := uuid.New()
		idMissing := uuid.New()

		mockInsight := database.IMInsightModel{Id: idExisting, Text: "Only Found Insight"}

		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{
				idExisting: mockInsight,
			},
		}

		repo := insight.NewInMemoryDBInsightRepository(db)

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{idExisting, idMissing})

		// Assert
		assert.Len(t, result, 1)
		assert.Equal(t, insight.Insight{Id: idExisting, Text: "Only Found Insight"}, result[0])
	})

	t.Run("should return empty list when no ids match", func(t *testing.T) {
		// Arrange
		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{},
		}
		repo := insight.NewInMemoryDBInsightRepository(db)
		missingID := uuid.New()

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{missingID})

		// Assert
		assert.Empty(t, result)
	})

	t.Run("should return empty list when ids slice is empty", func(t *testing.T) {
		// Arrange
		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{},
		}
		repo := insight.NewInMemoryDBInsightRepository(db)

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{})

		// Assert
		assert.Empty(t, result)
	})
}

func TestInMemoryDBInsightRepository_GetById(t *testing.T) {
	t.Run("should return insight when ID exists", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		mockInsight := database.IMInsightModel{Id: id, Text: "Existing Insight"}

		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{
				id: mockInsight,
			},
		}
		repo := insight.NewInMemoryDBInsightRepository(db)

		expected := &insight.Insight{Id: id, Text: "Existing Insight"}

		// Act
		result, err := repo.GetById(id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return error when ID does not exist", func(t *testing.T) {
		// Arrange
		db := &database.IMDatabase{
			InsightStorage: map[uuid.UUID]database.IMInsightModel{},
		}
		repo := insight.NewInMemoryDBInsightRepository(db)

		missingID := uuid.New()

		// Act
		result, err := repo.GetById(missingID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
