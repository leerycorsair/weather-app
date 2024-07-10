package postgres

import (
	"fmt"
	"weather-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type ForecastRepository struct {
	db *sqlx.DB
}

func NewForecastRepository(db *sqlx.DB) *ForecastRepository {
	return &ForecastRepository{db: db}
}

func (r *ForecastRepository) CreateForecast(forecast models.Forecast) (int, error) {
	var id int
	query := fmt.Sprintf(`
		insert into %s (city_id, temp, date, forecast_json)
		values ($1, $2, $3, $4)
		on conflict (city_id, date) do update set
			temp = excluded.temp,
			forecast_json = excluded.forecast_json
		returning id
	`, ForecastsTable)
	row := r.db.QueryRow(query, forecast.CityId, forecast.Temp, forecast.Date, forecast.ForecastJson)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ForecastRepository) GetForecasts(cityId int) ([]models.Forecast, error) {
	var forecasts []models.Forecast
	query := fmt.Sprintf("select id, city_id, temp, date, forecast_json from %s where city_id=$1", ForecastsTable)
	err := r.db.Select(&forecasts, query, cityId)
	if err != nil {
		return nil, err
	}
	return forecasts, nil
}
