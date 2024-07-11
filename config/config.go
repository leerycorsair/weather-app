package config

import (
	"flag"
	"os"
	"time"
	"weather-app/internal/repository/postgres"
)

func LoadOpenWeatherAPIKey() (string, error) {
	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	return apiKey, nil
}

type CollectorFlags struct {
	ToStart    bool
	Filename   string
	UpdateTime time.Duration
	Parallel   bool
}

func ParseCollectorFlags() (CollectorFlags, error) {
	var (
		toStart    bool
		filename   string
		updateTime string
		parallel   bool
	)
	flag.BoolVar(&toStart, "s", false, "Start data collector")
	flag.StringVar(&filename, "f", "", "File containing list of cities")
	flag.StringVar(&updateTime, "u", "1m", "Update interval")
	flag.BoolVar(&parallel, "p", false, "Enable parallel mode")
	flag.Parse()

	interval, err := time.ParseDuration(updateTime)
	if err != nil {
		return CollectorFlags{}, err
	}

	return CollectorFlags{
		ToStart:    toStart,
		Filename:   filename,
		UpdateTime: interval,
		Parallel:   parallel,
	}, nil
}

func LoadPGConfig() (postgres.PGConfig, error) {
	return postgres.PGConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PSWD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}, nil
}
