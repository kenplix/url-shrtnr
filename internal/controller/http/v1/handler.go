package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Kenplix/url-shrtnr/internal/service"
)

type Handler struct {
	services *service.Services
	logger   *zap.Logger
}

func NewHandler(logger *zap.Logger, services *service.Services) (*Handler, error) {
	if logger == nil {
		return nil, errors.New("logger not provided")
	}

	if services == nil {
		return nil, errors.New("services not provided")
	}

	h := &Handler{
		services: services,
		logger:   logger,
	}

	return h, nil
}

func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	h.initAuthRoutes(v1)
	h.initUsersRoutes(v1)
}
