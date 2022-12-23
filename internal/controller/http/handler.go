package http

import (
	"log"
	"net"

	"github.com/kenplix/url-shrtnr/internal/config"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/docs"
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

type Config struct {
	Host string
	Port string
}

func (h *Handler) InitEngine(env config.Environment, cfg Config) *gin.Engine {
	router := gin.New()

	host := cfg.Host
	if host == "" {
		host = "localhost"
	}

	docs.SwaggerInfo.Host = net.JoinHostPort(host, cfg.Port)
	log.Printf("shost: %s, host: %s", docs.SwaggerInfo.Host, host)

	if env == config.ProductionEnvironment {
		gin.SetMode(gin.ReleaseMode)
	} else {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

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
