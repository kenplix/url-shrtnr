package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// testGinContext returns gin context mock
func testGinContext(t *testing.T, w *httptest.ResponseRecorder) *gin.Context {
	t.Helper()

	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	ctx := log.ContextWithLogger(c.Request.Context(), testLogger(t))
	c.Request = c.Request.WithContext(ctx)

	return c
}

func testLogger(t *testing.T) *zap.Logger {
	t.Helper()

	logger, err := log.NewLogger(log.SetLevel(zap.DebugLevel.String()))
	require.NoErrorf(t, err, "failed to create testing logger: %s", err)

	return logger
}
