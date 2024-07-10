package repository

import (
	"weather-app/internal/models"
)

type CityRepository interface {
	CreateCity(models.City) (int, error)
	GetCities() ([]models.City, error)
	GetCity(cityId int) (models.City, error)
}

type ForecastRepository interface {
	CreateForecast(models.Forecast) (int, error)
	GetForecasts(cityId int) ([]models.Forecast, error)
}

type UserRepository interface {
	CreateUser(user models.User) (int, error)
	GetUser(login, password string) (models.User, error)
	GetFavorites(userId int) ([]int, error)
	AddFavorite(userId int, cityId int) (int, error)
	DeleteFavorite(userId int, cityId int) error
}

type Repository struct {
	CityRepository
	ForecastRepository
	UserRepository
}

func NewRepository(cityRep CityRepository, forecastRep ForecastRepository, userRep UserRepository) *Repository {
	return &Repository{
		CityRepository:     cityRep,
		ForecastRepository: forecastRep,
		UserRepository:     userRep,
	}
}
