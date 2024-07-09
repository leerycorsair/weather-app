package handler

import (
	"weather-app/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUpUser)
		auth.POST("/sign-in", h.signInUser)
	}

	api := router.Group("/api")
	{
		users := api.Group("/users", h.identifyUser)
		{
			users.GET("/favorites", h.getFavorites)
			users.POST("/favorites", h.addFavorite)
			users.DELETE("/favorites", h.deleteFavorite)
		}

		cities := api.Group("/cities")
		{
			cities.GET("", h.getCities)
		}

		forecasts := api.Group("/forecast")
		{
			forecasts.GET("/short/:city_id", h.getShortForecast)
			forecasts.GET("/detailed/:city_id", h.getDetailedForecast)
		}
	}

	return router
}
