package chart_test

import (
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/database"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetByIds(t *testing.T) {
	t.Run("should return all charts when all IDs exist", func(t *testing.T) {
		// Arrange
		chart1ID := uuid.New()
		chart2ID := uuid.New()

		mockDB := &database.InMemoryDatabase{
			ChartStorage: map[uuid.UUID]database.ChartModel{
				chart1ID: {
					Id:         chart1ID,
					Title:      "Chart 1",
					XAxisTitle: "X Axis 1",
					YAxisTitle: "Y Axis 1",
					Data: []map[string]float64{
						{"x": 1, "y": 1},
						{"x": 2, "y": 3},
					},
				},
				chart2ID: {
					Id:         chart2ID,
					Title:      "Chart 2",
					XAxisTitle: "X Axis 2",
					YAxisTitle: "Y Axis 2",
					Data: []map[string]float64{
						{"x": 1, "y": 1},
						{"x": 3, "y": 2},
					},
				},
			},
		}
		repo := chart.NewInMemoryDBChartRepository(mockDB)

		expectedResult := []chart.Chart{
			{
				Id:         chart1ID,
				Title:      "Chart 1",
				XAxisTitle: "X Axis 1",
				YAxisTitle: "Y Axis 1",
				Data: chart.ChartData{
					{"x": 1, "y": 1},
					{"x": 2, "y": 3},
				},
			},
			{
				Id:         chart2ID,
				Title:      "Chart 2",
				XAxisTitle: "X Axis 2",
				YAxisTitle: "Y Axis 2",
				Data: chart.ChartData{
					{"x": 1, "y": 1},
					{"x": 3, "y": 2},
				},
			},
		}

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{chart1ID, chart2ID})

		// Assert
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should return only matching charts when some IDs do not exist", func(t *testing.T) {
		// Arrange
		chart1ID := uuid.New()
		nonexistentID := uuid.New()

		mockDB := &database.InMemoryDatabase{
			ChartStorage: map[uuid.UUID]database.ChartModel{
				chart1ID: {
					Id:         chart1ID,
					Title:      "Chart 1",
					XAxisTitle: "X Axis 1",
					YAxisTitle: "Y Axis 1",
					Data: []map[string]float64{
						{"x": 1, "y": 1},
						{"x": 2, "y": 3},
					},
				},
			},
		}
		repo := chart.NewInMemoryDBChartRepository(mockDB)
		expectedResult := []chart.Chart{
			{
				Id:         chart1ID,
				Title:      "Chart 1",
				XAxisTitle: "X Axis 1",
				YAxisTitle: "Y Axis 1",
				Data: chart.ChartData{
					{"x": 1, "y": 1},
					{"x": 2, "y": 3},
				},
			},
		}

		// Act
		result, _ := repo.GetByIds(uuid.UUIDs{chart1ID, nonexistentID})

		// Assert
		assert.Equal(t, expectedResult, result)
	})

	t.Run("should return empty slice when input is empty", func(t *testing.T) {
		// Arrange
		mockDB := &database.InMemoryDatabase{
			ChartStorage: map[uuid.UUID]database.ChartModel{},
		}
		repo := chart.NewInMemoryDBChartRepository(mockDB)

		// Act
		result, _ := repo.GetByIds([]uuid.UUID{})

		// Assert
		assert.Empty(t, result)
	})
}
