package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"weather-app/config"
	datacollector "weather-app/internal/data_collector"
	"weather-app/internal/handler"
	"weather-app/internal/repository/postgres"
	"weather-app/internal/service"
	cityservice "weather-app/internal/service/city_service"
	forecastservice "weather-app/internal/service/forecast_service"
	userservice "weather-app/internal/service/user_service"
	"weather-app/server"

	_ "weather-app/docs"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// @title WebApp API
// @version 1.0
// @describtion API Server for WebApp

// @host localhost:8000
// BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	pgCfg, err := config.LoadPGConfig()
	if err != nil {
		logrus.Fatalf("Failed to load PG config: %v", err)
	}
	db, err := postgres.NewPgConnection(pgCfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Fatalf("Failed to shut down the database: %v", err)
		}
	}()

	cityRep := postgres.NewCityRepository(db)
	forecastRep := postgres.NewForecastRepository(db)
	userRep := postgres.NewUserRepository(db)

	cityServ := cityservice.NewCityService(cityRep)
	forecastServ := forecastservice.NewForecastService(cityServ, forecastRep)
	userServ := userservice.NewUserService(cityServ, userRep)
	service := service.NewService(userServ, cityServ, forecastServ)

	collectorCfg, err := config.ParseCollectorFlags()
	if err != nil {
		logrus.Fatalf("Failed to parse datacollector flags: %v", err)
	}

	if collectorCfg.ToStart {
		apiKey, err := config.LoadOpenWeatherAPIKey()
		if err != nil {
			logrus.Fatalf("Failed to load apiKey: %v", err)
		}
		dataCollector := datacollector.NewDataCollector(collectorCfg, service, apiKey)
		go dataCollector.Start()
	}

	handler := handler.NewHandler(service)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(os.Getenv("SERVER_PORT"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("ERROR running server:%s", err.Error())
		}
	}()
	defer func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Fatalf("ERROR shutting down server: %s", err.Error())
		}
	}()

	logrus.Print("WebApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("WebApp Shutting Down")
}
