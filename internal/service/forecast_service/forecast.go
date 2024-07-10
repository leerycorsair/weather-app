package forecastservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
	"weather-app/internal/models"
	"weather-app/internal/repository"
	"weather-app/internal/service"
)

type ForecastService struct {
	cityService service.CityService
	forecastRep repository.ForecastRepository
}

func NewForecastService(cityservice service.CityService, forecastRep repository.ForecastRepository) *ForecastService {
	return &ForecastService{
		cityService: cityservice,
		forecastRep: forecastRep,
	}
}

func (s *ForecastService) CreateForecast(forecast models.Forecast) (int, error) {
	return s.forecastRep.CreateForecast(forecast)
}

func filterFutureForecasts(forecasts []models.Forecast) []models.Forecast {
	now := time.Now()
	var futureForecasts []models.Forecast
	for _, f := range forecasts {
		if f.Date.After(now) {
			futureForecasts = append(futureForecasts, f)
		}
	}
	return futureForecasts
}

func (s *ForecastService) GetShortForecast(cityId int) (models.ForecastSummary, error) {
	var summary models.ForecastSummary
	city, err := s.cityService.GetCity(cityId)
	if err != nil {
		return summary, err
	}
	summary.City, summary.Country = city.Name, city.Country
	forecasts, err := s.forecastRep.GetForecasts(cityId)
	if err != nil {
		return summary, err
	}
	forecasts = filterFutureForecasts(forecasts)

	dates := make(map[time.Time]struct{})
	for _, f := range forecasts {
		date := time.Date(f.Date.Year(), f.Date.Month(), f.Date.Day(), 0, 0, 0, 0, time.UTC)
		dates[date] = struct{}{}
	}
	for date := range dates {
		summary.AvailableDates = append(summary.AvailableDates, date.Format("2006-01-02"))
	}
	sort.Slice(summary.AvailableDates, func(i, j int) bool {
		return summary.AvailableDates[i] < summary.AvailableDates[j]
	})

	for _, f := range forecasts {
		summary.AvgTemp += f.Temp
	}

	if len(forecasts) > 0 {
		summary.AvgTemp /= float32(len(forecasts))
	} else {
		summary.AvgTemp = 0
	}
	return summary, nil
}

func filterForecastsByDateTime(forecasts []models.Forecast, target time.Time) []models.Forecast {
	var filtered []models.Forecast
	hasTimeComponent := !(target.Hour() == 0 && target.Minute() == 0 && target.Second() == 0)
	for _, forecast := range forecasts {
		if forecast.Date.Year() == target.Year() &&
			forecast.Date.Month() == target.Month() &&
			forecast.Date.Day() == target.Day() {
			if hasTimeComponent {
				if forecast.Date.Hour() == target.Hour() &&
					forecast.Date.Minute() == target.Minute() &&
					forecast.Date.Second() == target.Second() {
					filtered = append(filtered, forecast)
				}
			} else {
				filtered = append(filtered, forecast)
			}
		}
	}
	return filtered
}

func (s *ForecastService) GetDetailedForecast(cityId int, date time.Time) ([]models.Forecast, error) {
	forecasts, err := s.forecastRep.GetForecasts(cityId)
	if err != nil {
		return nil, err
	}
	filtered := filterForecastsByDateTime(forecasts, date)
	if len(filtered) == 0 {
		return nil, errors.New("no forecasts were found")
	}
	sort.Slice(forecasts, func(i, j int) bool {
		return forecasts[i].Date.Before(forecasts[j].Date)
	})
	return filtered, nil
}

type weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type mainData struct {
	Temp      float32 `json:"temp"`
	FeelsLike float32 `json:"feels_like"`
	TempMin   float32 `json:"temp_min"`
	TempMax   float32 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type forecastItem struct {
	Dt      int64     `json:"dt"`
	Main    mainData  `json:"main"`
	Weather []weather `json:"weather"`
	Clouds  struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float32 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float32 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float32 `json:"pop"`
	Rain       struct {
		ThreeH float32 `json:"3h,omitempty"`
	} `json:"rain"`
	Sys struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
}

type forecastResponse struct {
	List []forecastItem `json:"list"`
}

func (s *ForecastService) FetchForecastData(city models.City, openWeatherAPIKey string) ([]models.Forecast, error) {
	var forecasts []models.Forecast

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&units=metric&appid=%s", city.Latitude, city.Longitude, openWeatherAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to OpenWeather API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var forecastResponse forecastResponse
	if err := json.Unmarshal(body, &forecastResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	for _, item := range forecastResponse.List {
		date := time.Unix(item.Dt, 0)
		forecastJson, _ := json.Marshal(item)

		forecast := models.Forecast{
			CityId:       city.Id,
			Temp:         item.Main.Temp,
			Date:         date,
			ForecastJson: forecastJson,
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}
