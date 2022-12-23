package http

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kenplix/url-shrtnr/internal/controller/http/ginctx"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kenplix/url-shrtnr/pkg/log"
)

func requestIDMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return requestid.New(
		requestid.WithHandler(func(c *gin.Context, requestID string) {
			withRequestID := logger.With(zap.String("request-id", requestID))

			ctx := log.ContextWithLogger(c.Request.Context(), withRequestID)
			c.Request = c.Request.WithContext(ctx)
		}),
	)
}

type requestReader struct {
	io.ReadCloser
	buf *bytes.Buffer
}

func (r *requestReader) Read(p []byte) (n int, err error) {
	return io.TeeReader(r.ReadCloser, r.buf).Read(p)
}

func requestReaderMiddleware(c *gin.Context) {
	c.Request.Body = &requestReader{
		ReadCloser: c.Request.Body,
		buf:        &bytes.Buffer{},
	}
}

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return io.MultiWriter(w.buf, w.ResponseWriter).Write(p)
}

func responseWriterMiddleware(c *gin.Context) {
	c.Writer = &responseWriter{
		ResponseWriter: c.Writer,
		buf:            &bytes.Buffer{},
	}
}

func loggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		SkipPaths: []string{
			"/api/v1/auth/sign-up",
			"/api/v1/auth/sign-in",
			"/api/v1/auth/refresh-tokens",
		},
		Context: func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			if requestID := requestid.Get(c); requestID != "" {
				fields = append(fields, zap.String("request-id", requestID))
			}

			r := c.Request.Body.(*requestReader)
			fields = append(fields, zap.String("request-body", r.buf.String()))

			w := c.Writer.(*responseWriter)
			fields = append(fields, zap.String("response-body", w.buf.String()))

			return fields
		},
	})
}

func translatorMiddleware(unitrans *ut.UniversalTranslator) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := c.Query("locale")
		languages := parseAcceptLanguageHeader(c)

		translator, _ := unitrans.FindTranslator(append([]string{locale}, languages...)...)
		c.Set(ginctx.TranslatorContext, translator)

		logger := log.LoggerFromContext(c.Request.Context())

		logger.Debug("user locale successfully set",
			zap.String("locale", translator.Locale()),
		)
	}
}

// parseAcceptLanguageHeader returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func parseAcceptLanguageHeader(c *gin.Context) []string {
	header := c.GetHeader("Accept-Language")
	if header == "" {
		return nil
	}

	options := strings.Split(header, ",")
	languages := make([]string, 0, len(options))

	for _, option := range options {
		locale := strings.SplitN(option, ";", 2)
		languages = append(languages, strings.Trim(locale[0], " "))
	}

	return languages
}

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
