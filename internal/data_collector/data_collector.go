package datacollector

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"weather-app/config"
	"weather-app/internal/models"
	"weather-app/internal/service"

	"github.com/sirupsen/logrus"
)

type DataCollector struct {
	services   *service.Service
	citiesFile string
	updateTime time.Duration
	parallel   bool
	apiKey     string
}

func NewDataCollector(cfg config.CollectorFlags, services *service.Service, apiKey string) *DataCollector {
	return &DataCollector{
		services:   services,
		citiesFile: cfg.Filename,
		updateTime: cfg.UpdateTime,
		parallel:   cfg.Parallel,
		apiKey:     apiKey,
	}
}

func (dc *DataCollector) Start() {
	if dc.citiesFile != "" {
		citiesNames, err := readLines(dc.citiesFile)
		if err != nil {
			logrus.Errorf("Failed to load cities from file: %v", err)
		}
		dc.fetchAndCreateCities(citiesNames)
		logrus.Printf("Cities data was updated at %v", time.Now())
	}

	cities, err := dc.services.CityService.GetCities()
	if err != nil {
		logrus.Fatalf("Failed to load cities from database: %v", err)
	}
	if len(cities) == 0 {
		logrus.Fatalf("No cities found in the database")
	}

	dc.fetchAndCreateForecasts(cities)
	logrus.Printf("Weather was updated at %v", time.Now())

	ticker := time.NewTicker(dc.updateTime)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			dc.fetchAndCreateForecasts(cities)
			logrus.Printf("Weather was updated at %v", time.Now())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

func (dc *DataCollector) fetchAndCreateCities(citiesNames []string) {
	if dc.parallel {
		var wg sync.WaitGroup
		errorChan := make(chan error, len(citiesNames))

		for _, cityName := range citiesNames {
			wg.Add(1)
			go func(cityName string) {
				defer wg.Done()
				city, err := dc.services.CityService.FetchCityData(cityName, dc.apiKey)
				if err != nil {
					errorChan <- fmt.Errorf("Failed to fetch city data for %v: %v", cityName, err)
					return
				}
				_, err = dc.services.CityService.CreateCity(city)
				if err != nil {
					errorChan <- fmt.Errorf("Failed to create city record in db: %v", err)
					return
				}
			}(cityName)
		}

		wg.Wait()
		close(errorChan)

		for err := range errorChan {
			logrus.Error(err)
		}
	} else {
		for _, cityName := range citiesNames {
			city, err := dc.services.CityService.FetchCityData(cityName, dc.apiKey)
			if err != nil {
				logrus.Errorf("Failed to fetch city data for %v: %v", cityName, err)
				continue
			}
			_, err = dc.services.CityService.CreateCity(city)
			if err != nil {
				logrus.Fatalf("Failed to create city record in db: %v", err)
			}
		}
	}
}

func (dc *DataCollector) fetchAndCreateForecasts(cities []models.City) {
	if dc.parallel {
		var wg sync.WaitGroup
		errorChan := make(chan error, len(cities))

		for _, city := range cities {
			wg.Add(1)
			go func(city models.City) {
				defer wg.Done()
				forecasts, err := dc.services.ForecastService.FetchForecastData(city, dc.apiKey)
				if err != nil {
					errorChan <- fmt.Errorf("Failed to fetch forecast data for %v: %v", city.Name, err)
					return
				}
				for _, forecast := range forecasts {
					_, err := dc.services.ForecastService.CreateForecast(forecast)
					if err != nil {
						errorChan <- fmt.Errorf("Failed to create forecast record in db: %v", err)
						return
					}
				}
			}(city)
		}

		wg.Wait()
		close(errorChan)

		for err := range errorChan {
			logrus.Error(err)
		}
	} else {
		for _, city := range cities {
			forecasts, err := dc.services.ForecastService.FetchForecastData(city, dc.apiKey)
			if err != nil {
				logrus.Errorf("Failed to fetch forecast data for %v: %v", city.Name, err)
				continue
			}
			for _, forecast := range forecasts {
				_, err := dc.services.ForecastService.CreateForecast(forecast)
				if err != nil {
					logrus.Fatalf("Failed to create forecast record in db: %v", err)
				}
			}
		}
	}
}

func readLines(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}
