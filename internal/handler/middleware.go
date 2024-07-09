package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) identifyUser(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "Empty auth header")
		return
	}

	// headerParts := strings.Split(header, " ")
	// if len(headerParts) != 2 {
	// 	logrus.Println(header)
	// 	logrus.Println(len(headerParts))
	// 	newErrorResponse(c, http.StatusUnauthorized, "Invalid auth header")
	// 	return
	// }
	userId, err := h.services.UserService.ParseToken(header)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set(userCtx, userId)
}
