package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Kenplix/url-shrtnr/internal/usecase"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
)

func NewHandler(manager *usecase.Manager, tokenService auth.TokenService) *gin.Engine {
	router := gin.New()

	v1 := router.Group("/v1")
	NewUsersHandler(manager.Users, tokenService).initRoutes(v1)

	return router
}
