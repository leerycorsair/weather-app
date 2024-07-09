package service

import (
	"time"
	"weather-app/internal/models"
)

type UserService interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(login string, password string) (string, error)
	ParseToken(accessToken string) (int, error)
	GetFavorites(userId int) ([]int, error)
	AddFavorite(userId int, cityId int) (int, error)
	DeleteFavorite(userId int, cityId int) error
}

type CityService interface {
	CreateCity(models.City) (int, error)
	GetCities() ([]models.City, error)
	GetCity(cityId int) (models.City, error)
	FetchCityData(cityName string, openWeatherAPIKey string) (models.City, error)
}

type ForecastService interface {
	CreateForecast(models.Forecast) (int, error)
	GetShortForecast(cityId int) (models.ForecastSummary, error)
	GetDetailedForecast(cityId int, date time.Time) ([]models.Forecast, error)
	FetchForecastData(city models.City, openWeatherAPIKey string) ([]models.Forecast, error)
}

type Service struct {
	UserService
	CityService
	ForecastService
}

func NewService(userService UserService, cityService CityService, forecastService ForecastService) *Service {
	return &Service{
		UserService:     userService,
		CityService:     cityService,
		ForecastService: forecastService,
	}
}
