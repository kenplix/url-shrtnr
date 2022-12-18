package http

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/internal/controller/http/validator"

	v1 "github.com/kenplix/url-shrtnr/internal/controller/http/v1"
	"github.com/kenplix/url-shrtnr/internal/service"
)

type Handler struct {
	v1       *v1.Handler
	unitrans *ut.UniversalTranslator
	logger   *zap.Logger
}

func NewHandler(logger *zap.Logger, services *service.Services) (*Handler, error) {
	if logger == nil {
		return nil, errors.New("logger not provided")
	}

	unitrans, err := validator.Init(logger)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to configure gin validator instance")
	}

	handlerV1, err := v1.NewHandler(logger, services)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API v1 handler")
	}

	h := &Handler{
		v1:       handlerV1,
		unitrans: unitrans,
		logger:   logger,
	}

	return h, nil
}

func (h *Handler) InitEngine() *gin.Engine {
	router := gin.New()

	router.Use(
		requestIDMiddleware(h.logger),
		requestReaderMiddleware,
		responseWriterMiddleware,
		loggerMiddleware(h.logger),
		corsMiddleware(),
		translatorMiddleware(h.unitrans),
	)

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	api := router.Group("/api")

	h.v1.InitRoutes(api)
}
