package v1

import (
	"net/http"

	"github.com/kenplix/url-shrtnr/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(router *gin.RouterGroup) {
	users := router.Group(
		"/users",
		h.userIdentityMiddleware,
		h.userActivityMiddleware,
	)

	users.GET("/me", h.me)
}

func (h *Handler) me(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	c.JSON(http.StatusOK, user)
}
