package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Kenplix/url-shrtnr/internal/usecase"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
)

func NewHandler(manager *usecase.Manager, tokensServ auth.TokensService) *gin.Engine {
	router := gin.New()

	v1 := router.Group("/v1")
	NewUsersHandler(manager.Users, tokensServ).initRoutes(v1)

	return router
}
