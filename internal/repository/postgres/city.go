package postgres

import (
	"fmt"
	"weather-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type CityRepository struct {
	db *sqlx.DB
}

func NewCityRepository(db *sqlx.DB) *CityRepository {
	return &CityRepository{db: db}
}

func (r *CityRepository) CreateCity(city models.City) (int, error) {
	var id int
	query := fmt.Sprintf(`
		insert into %s (name, country, latitude, longitude)
		values ($1, $2, $3, $4)
		on conflict (name, country) do update set
			latitude = excluded.latitude,
			longitude = excluded.longitude
		returning id
	`, CitiesTable)
	row := r.db.QueryRow(query, city.Name, city.Country, city.Latitude, city.Longitude)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *CityRepository) GetCities() ([]models.City, error) {
	var cities []models.City
	query := fmt.Sprintf("select id, name, country, latitude, longitude from %s order by name", CitiesTable)
	err := r.db.Select(&cities, query)
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *CityRepository) GetCity(cityId int) (models.City, error) {
	var city models.City
	query := fmt.Sprintf("select id, name, country, latitude, longitude from %s where id=$1", CitiesTable)
	err := r.db.Get(&city, query, cityId)
	if err != nil {
		return models.City{}, err
	}
	return city, nil
}
