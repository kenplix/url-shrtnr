package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/Kenplix/url-shrtnr/internal/usecase"
)

func NewHandler(manager *usecase.Manager) *gin.Engine {
	router := gin.New()

	v1 := router.Group("/v1")
	NewUsersHandler(manager.Users).initRoutes(v1)

	return router
}
