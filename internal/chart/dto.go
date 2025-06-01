package chart

import "github.com/google/uuid"

type (
	ChartData []map[string]float64
	Chart     struct {
		Id         uuid.UUID            `json:"id"`
		Title      string               `json:"title"`
		XAxisTitle string               `json:"x_axis_title"`
		YAxisTitle string               `json:"y_axis_title"`
		Data       []map[string]float64 `json:"data"`
	}
)
