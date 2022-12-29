package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kenplix/url-shrtnr/internal/controller/http/validator"

	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func initValidator(t *testing.T) {
	t.Helper()

	_, err := validator.Init(testLogger(t))
	require.NoErrorf(t, err, "failed to initialize validator: %s", err)
}

func testLoggerMiddleware(t *testing.T) gin.HandlerFunc {
	t.Helper()

	return func(c *gin.Context) {
		ctx := log.ContextWithLogger(c.Request.Context(), testLogger(t))
		c.Request = c.Request.WithContext(ctx)
	}
}

func testLogger(t *testing.T) *zap.Logger {
	t.Helper()

	logger, err := log.NewLogger(log.SetLevel(zap.DebugLevel.String()))
	require.NoErrorf(t, err, "failed to create testing logger: %s", err)

	return logger
}

func testUnauthorizedErrorResponse(t *testing.T) string {
	t.Helper()

	return mustMarshal(t, errResponse{
		Errors: []apiError{newUnauthorizedError()},
	})
}

func testSuspendedErrorResponse(t *testing.T) string {
	t.Helper()

	return mustMarshal(t, errResponse{
		Errors: []apiError{newSuspendedError()},
	})
}

func testInternalErrorResponse(t *testing.T) string {
	t.Helper()

	return mustMarshal(t, errResponse{
		Errors: []apiError{newInternalError()},
	})
}

func mustMarshal(t *testing.T, data any) string {
	t.Helper()

	buf, err := json.Marshal(data)
	require.NoErrorf(t, err, "failed to marshal %v data", err)

	return string(buf)
}
