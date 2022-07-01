package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) init() {
	h.router.GET("/ping", func(c *gin.Context) {
		h.log.Debug("ping")
		c.String(http.StatusOK, "pong")
	})
}
