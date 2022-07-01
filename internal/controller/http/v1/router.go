package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) init() {
	h.router.GET("/ping", func(c *gin.Context) {
		h.log.Debugf("ping")
		c.String(http.StatusOK, "pong")
	})
}
