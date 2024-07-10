package cityservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weather-app/internal/models"
	"weather-app/internal/repository"
)

type CityService struct {
	cityRep repository.CityRepository
}

func NewCityService(cityRep repository.CityRepository) *CityService {
	return &CityService{
		cityRep: cityRep,
	}
}

func (s *CityService) CreateCity(city models.City) (int, error) {
	return s.cityRep.CreateCity(city)
}

func (s *CityService) GetCities() ([]models.City, error) {
	return s.cityRep.GetCities()
}
func (s *CityService) GetCity(cityId int) (models.City, error) {
	return s.cityRep.GetCity(cityId)
}

type geocodingResponse struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func (s *CityService) FetchCityData(cityName string, openWeatherAPIKey string) (models.City, error) {
	var city models.City
	url := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", cityName, openWeatherAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return city, fmt.Errorf("failed to make request to OpenWeather Geocoding API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return city, fmt.Errorf("failed to read response body: %w", err)
	}

	var geocodingResponses []geocodingResponse
	if err := json.Unmarshal(body, &geocodingResponses); err != nil {
		return city, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(geocodingResponses) == 0 {
		return city, fmt.Errorf("no results found for city: %s", cityName)
	}

	geo := geocodingResponses[0]
	city = models.City{
		Name:      geo.Name,
		Country:   geo.Country,
		Latitude:  geo.Lat,
		Longitude: geo.Lon,
	}

	return city, nil
}
