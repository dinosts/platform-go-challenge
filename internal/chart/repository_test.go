package chart_test

import (
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/database"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryDBChartRepository_GetByIds(t *testing.T) {
	chart1ID := uuid.New()
	chart2ID := uuid.New()
	chart3ID := uuid.New()

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

	tests := []struct {
		name           string
		inputIDs       []uuid.UUID
		expectedCharts []chart.Chart
	}{
		{
			name:     "should return charts when IDs exist",
			inputIDs: []uuid.UUID{chart1ID, chart2ID},
			expectedCharts: []chart.Chart{
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
			},
		},
		{
			name:     "should return partial charts when an chart was not found",
			inputIDs: []uuid.UUID{chart1ID, chart3ID},
			expectedCharts: []chart.Chart{
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
			},
		},
		{
			name:           "should return empty slice when input is empty",
			inputIDs:       []uuid.UUID{},
			expectedCharts: []chart.Chart{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			actualCharts := repo.GetByIds(tt.inputIDs)

			// Assert
			assert.Equal(t, tt.expectedCharts, actualCharts)
		})
	}
}
