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

// @Summary		Get users personal information
// @Security		JWT-RS256
// @Tags			user
// @Description	Get users personal information
// @Produce		json
// @Success		200	{object}	entity.User
// @Failure		401	{object}	errResponse{errors=[]entity.CoreError}	"Access is denied due to invalid credentials"
// @Failure		403	{object}	errResponse{errors=[]entity.CoreError}	"Your account has been suspended"
// @Failure		500	{object}	errResponse{errors=[]entity.CoreError}	"Internal server error"
// @Router			/users/me [get]
func (h *Handler) me(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	c.JSON(http.StatusOK, user)
}
