package v1

import (
	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
)

type UsersHandler struct {
	usersServ  service.UsersService
	tokensServ auth.TokensService
}

func NewUsersHandler(usersServ service.UsersService, tokensServ auth.TokensService) (*UsersHandler, error) {
	if usersServ == nil {
		return nil, errors.New("users service not provided")
	}

	if tokensServ == nil {
		return nil, errors.New("tokens service not provided")
	}

	h := &UsersHandler{
		usersServ:  usersServ,
		tokensServ: tokensServ,
	}

	return h, nil
}

func (h *UsersHandler) init(router *gin.RouterGroup) {
	usersGroup := router.Group("/users", h.identityMiddleware)
	{
		usersGroup.GET("/me", h.me)
	}
}

func (h *UsersHandler) me(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		log.Printf("error: failed to get user profile: %s", err)
		errorResponse(c, http.StatusInternalServerError, &entity.CoreError{
			Code:    errorcode.InternalError,
			Message: http.StatusText(http.StatusInternalServerError),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userID": id.Hex(),
	})
}

func getUserID(c *gin.Context) (primitive.ObjectID, error) {
	return getIDByContext(c, userContext)
}

func (h *UsersHandler) identityMiddleware(c *gin.Context) {
	userID, err := h.parseAuthorizationHeader(c)
	if err != nil {
		log.Printf(`warning: error parsing "Authorization" header: %s`, err)
		errorResponse(c, http.StatusUnauthorized, &entity.CoreError{
			Code:    errorcode.UnauthorizedAccess,
			Message: "access is denied due to invalid credentials",
		})

		return
	}

	c.Set(userContext, userID)
}

func (h *UsersHandler) parseAuthorizationHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if headerParts[1] == "" {
		return "", errors.New("token is empty")
	}

	return h.tokensServ.ParseAccessToken(headerParts[1])
}
