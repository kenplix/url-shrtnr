package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/service"
)

type UsersHandler struct {
	usersServ service.UsersService
	jwtServ   service.JWTService
}

func NewUsersHandler(usersServ service.UsersService, jwtServ service.JWTService) (*UsersHandler, error) {
	if usersServ == nil {
		return nil, errors.New("users service not provided")
	}

	if jwtServ == nil {
		return nil, errors.New("jwt service not provided")
	}

	h := &UsersHandler{
		usersServ: usersServ,
		jwtServ:   jwtServ,
	}

	return h, nil
}

func (h *UsersHandler) init(router *gin.RouterGroup) {
	usersGroup := router.Group("/users", userIdentityMiddleware(h.jwtServ))

	usersGroup.GET("/me", h.me)
}

func (h *UsersHandler) me(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		log.Printf("error: failed to get userID object: %s", err)
		internalErrorResponse(c)

		return
	}

	user, err := h.usersServ.GetByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("error: failed to get user[id:%q]: %s", userID.Hex(), err)
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, user)
}
