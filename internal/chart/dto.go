package chart

import "github.com/google/uuid"

type (
	ChartData []map[string]float64
	Chart     struct {
		Id         uuid.UUID
		Title      string
		XAxisTitle string
		YAxisTitle string
		Data       []map[string]float64
	}
)
