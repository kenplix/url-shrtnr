package v1

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

const (
	userContext       = "userID"
	translatorContext = "localeTranslator"
)

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
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

func translatorMiddleware(c *gin.Context) {
	locale := c.Query("locale")
	languages := parseAcceptLanguageHeader(c)

	trans, _ := universalTranslator.FindTranslator(append([]string{locale}, languages...)...)
	c.Set(translatorContext, trans)
}

// parseAcceptLanguageHeader returns an array of accepted languages denoted by
// the Accept-Language header sent by the browser
func parseAcceptLanguageHeader(c *gin.Context) (languages []string) {
	header := c.GetHeader("Accept-Language")
	if header == "" {
		return
	}

	options := strings.Split(header, ",")
	languages = make([]string, 0, len(options))
	for _, option := range options {
		locale := strings.SplitN(option, ";", 2)
		languages = append(languages, strings.Trim(locale[0], " "))
	}

	return
}

func getIDByContext(c *gin.Context, context string) (primitive.ObjectID, error) {
	idStr := c.GetString(context)
	if idStr == "" {
		return primitive.NilObjectID, fmt.Errorf("%s from context is empty", context)
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, errors.Wrapf(err, "failed to get %s from hex", context)
	}

	return id, nil
}
