package models

import (
	"encoding/json"
	"time"
)

// Forecast represents the weather forecast model
// @Description Weather forecast model
type Forecast struct {
	Id           int             `json:"id"  db:"id"`                       // @Description Forecast ID
	CityId       int             `json:"city_id"  db:"city_id"`             // @Description City ID
	Temp         float32         `json:"temp"  db:"temp"`                   // @Description Temperature
	Date         time.Time       `json:"date"  db:"date"`                   // @Description Date of the forecast
	ForecastJson json.RawMessage `json:"forecast_json"  db:"forecast_json"` // @Description Raw JSON of the forecast details
}

// ForecastSummary represents a summary of weather forecasts
// @Description Weather forecast summary
type ForecastSummary struct {
	Country        string   `json:"country"  db:"country"`                 // @Description Country
	City           string   `json:"city"  db:"city"`                       // @Description City
	AvgTemp        float32  `json:"avg_temp"  db:"avg_temp"`               // @Description Average Temperature
	AvailableDates []string `json:"available_dates"  db:"available_dates"` // @Description Available dates for forecasts
}
