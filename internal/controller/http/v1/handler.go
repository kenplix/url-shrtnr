package v1

import (
	"net/http"

	"github.com/Kenplix/url-shrtnr/pkg/logger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	router *gin.Engine
	log    logger.Interface
}

func NewHandler(log logger.Interface) *Handler {
	h := &Handler{
		router: gin.New(),
		log:    log,
	}

	h.init()
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
