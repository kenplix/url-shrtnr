package v1

import (
	"log"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Kenplix/url-shrtnr/internal/service"
)

type Handler struct {
	authHandler  *AuthHandler
	usersHandler *UsersHandler
}

var universalTranslator *ut.UniversalTranslator

func init() {
	english := en.New()
	russian := ru.New()

	universalTranslator = ut.New(english, english, russian)

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Printf("info: configuring gin validator instance")

		if err := configureValidator(validate, universalTranslator); err != nil {
			panic(errors.Wrap(err, "failed to configure gin validator instance"))
		}
	}
}

func NewHandler(services *service.Services) (*Handler, error) {
	if services == nil {
		return nil, errors.New("services not provided")
	}

	authHandler, err := NewAuthHandler(services.Auth, services.Users, services.JWT)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth handler")
	}

	usersHandler, err := NewUsersHandler(services.Users, services.JWT)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create users handler")
	}

	h := &Handler{
		authHandler:  authHandler,
		usersHandler: usersHandler,
	}

	return h, nil
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	logger, _ := zap.NewDevelopment()

	router.Use(
		requestReaderMiddleware,
		responseWriterMiddleware,
		loggerMiddleware(logger),
		corsMiddleware(),
		translatorMiddleware,
	)

	handlers := []interface{ init(group *gin.RouterGroup) }{
		h.authHandler,
		h.usersHandler,
	}

	v1 := router.Group("api/v1")
	for _, handler := range handlers {
		handler.init(v1)
	}

	return router
}
