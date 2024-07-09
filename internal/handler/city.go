package handler

import (
	"net/http"
	"weather-app/internal/models"

	"github.com/gin-gonic/gin"
)

type GetCitiesResponse struct {
	Cities []models.City `json:"cities"  db:"cities"`
}

// getCities retrieves the list of cities
// @Summary Get cities
// @Description Get the list of cities
// @Tags cities
// @Produce json
// @Success 200 {object} GetCitiesResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/cities [get]
func (h *Handler) getCities(c *gin.Context) {
	cities, err := h.services.CityService.GetCities()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, GetCitiesResponse{
		Cities: cities,
	})
}
