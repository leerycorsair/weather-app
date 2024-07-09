package handler

import (
	"net/http"
	"strconv"
	"time"
	"weather-app/internal/models"

	"github.com/gin-gonic/gin"
)

type GetShortForecastResponse struct {
	Forecast models.ForecastSummary `json:"forecast"  db:"forecast"`
}

// getShortForecast retrieves the short forecast for a city
// @Summary Get short forecast
// @Description Get the short forecast for a specific city
// @Tags forecast
// @Produce json
// @Param city_id path int true "City ID"
// @Success 200 {object} GetShortForecastResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/forecast/short/{city_id} [get]
func (h *Handler) getShortForecast(c *gin.Context) {
	cityId, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	forecast, err := h.services.ForecastService.GetShortForecast(int(cityId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, GetShortForecastResponse{
		Forecast: forecast,
	})
}

type GetDetailedForecastResponse struct {
	City      string            `json:"city"  db:"city"`
	Forecasts []models.Forecast `json:"forecasts"  db:"forecasts"`
}

// getDetailedForecast retrieves the detailed forecast for a city on a specific date
// @Summary Get detailed forecast
// @Description Get the detailed forecast for a specific city on a specific date
// @Tags forecast
// @Produce json
// @Param city_id path int true "City ID"
// @Param date query string true "Date" Format(date)
// @Success 200 {object} GetDetailedForecastResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/forecast/detailed/{city_id} [get]
func (h *Handler) getDetailedForecast(c *gin.Context) {
	cityId, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	dateStr := c.Query("date")
	var date time.Time
	layouts := []string{"2006-01-02", "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		date, err = time.Parse(layout, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid date format. Use '2006-01-02' or '2006-01-02 15:04:05'")
		return
	}
	forecasts, err := h.services.ForecastService.GetDetailedForecast(int(cityId), date)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	city, err := h.services.CityService.GetCity(int(cityId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, GetDetailedForecastResponse{
		City:      city.Name,
		Forecasts: forecasts,
	})
}
