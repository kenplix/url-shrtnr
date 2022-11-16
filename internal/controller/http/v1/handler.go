package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
)

func NewHandler(services *service.Services, tokensServ auth.TokensService) *gin.Engine {
	router := gin.New()

	v1 := router.Group("/v1")
	NewUsersHandler(services.Users, tokensServ).initRoutes(v1)

	return router
}
