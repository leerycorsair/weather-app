package handler

import (
	"net/http"
	"strconv"
	"weather-app/internal/dto"
	"weather-app/internal/models"

	"github.com/gin-gonic/gin"
)

type SignUpUserResponse struct {
	Id int `json:"id"  db:"id"`
}

// signUpUser registers a new user
// @Summary Sign up user
// @Description Registers a new user with login, password, and email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.DTOSignUp true "Sign up info"
// @Success 200 {object} SignUpUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUpUser(c *gin.Context) {
	var user dto.DTOSignUp
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.UserService.CreateUser(models.User{
		Login:    user.Login,
		Password: user.Password,
		Email:    user.Email,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, SignUpUserResponse{id})
}

type SignInUserResponse struct {
	Token string `json:"token"  db:"token"`
}

// signInUser authenticates a user and returns a token
// @Summary Sign in user
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.DTOSignIn true "Sign in info"
// @Success 200 {object} SignInUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signInUser(c *gin.Context) {
	var user dto.DTOSignIn
	if err := c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.UserService.GenerateToken(user.Login, user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, SignInUserResponse{token})
}

type GetFavoritesResponse struct {
	Cities []models.City `json:"cities"  db:"cities"`
}

// getFavorites retrieves the list of favorite cities for the authenticated user
// @Summary Get favorite cities
// @Description Retrieves the list of favorite cities for the authenticated user
// @Tags favorites
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} GetFavoritesResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/favorites [get]
func (h *Handler) getFavorites(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "UserId not found")
		return
	}
	citiesIds, err := h.services.UserService.GetFavorites(userId.(int))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var resp GetFavoritesResponse
	for _, cityId := range citiesIds {
		city, err := h.services.CityService.GetCity(cityId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		resp.Cities = append(resp.Cities, city)
	}

	c.JSON(http.StatusOK, resp)
}

// addFavorite adds a city to the user's list of favorite cities
// @Summary Add favorite city
// @Description Adds a city to the user's list of favorite cities
// @Tags favorites
// @Produce json
// @Security ApiKeyAuth
// @Param cityId query int true "City ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/favorites [post]
func (h *Handler) addFavorite(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "UserId not found")
		return
	}
	cityId, err := strconv.ParseInt(c.Query("cityId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	_, err = h.services.UserService.AddFavorite(userId.(int), int(cityId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// deleteFavorite removes a city from the user's list of favorite cities
// @Summary Delete favorite city
// @Description Removes a city from the user's list of favorite cities
// @Tags favorites
// @Produce json
// @Security ApiKeyAuth
// @Param cityId query int true "City ID"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/users/favorites [delete]
func (h *Handler) deleteFavorite(c *gin.Context) {
	userId, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "UserId not found")
		return
	}
	cityId, err := strconv.ParseInt(c.Query("cityId"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UserService.DeleteFavorite(userId.(int), int(cityId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
