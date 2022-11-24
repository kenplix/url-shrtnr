package v1

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Kenplix/url-shrtnr/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	userContext       = "user"
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

func userIdentityMiddleware(usersServ service.UsersService, jwtServ service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := parseAuthorizationHeader(c)
		if err != nil {
			log.Printf(`warning: failed to parse "Authorization" header: %s`, err)
			unauthorizedErrorResponse(c)

			return
		}

		claims, err := jwtServ.ParseAccessToken(accessToken)
		if err != nil {
			log.Printf(`warning: failed to parse %q access token: %s`, accessToken, err)
			unauthorizedErrorResponse(c)

			return
		}

		ctx := c.Request.Context()

		err = jwtServ.ValidateAccessToken(ctx, claims)
		if err != nil {
			log.Printf("warning: failed to validate %+v access token: %s", claims, err)
			unauthorizedErrorResponse(c)

			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.Subject)
		if err != nil {
			log.Printf("warning: failed to get userID object from %q hex", claims.Subject)
			unauthorizedErrorResponse(c)

			return
		}

		user, err := usersServ.GetByID(ctx, userID)
		if err != nil {
			log.Printf("warning: failed to get user[id:%q]", userID.Hex())
			unauthorizedErrorResponse(c)

			return
		} else if user.SuspendedAt != nil {
			log.Printf("warning: protected route request from suspended user[id:%q]", userID.Hex())
			suspendedErrorResponse(c)

			return
		}

		jwtServ.ProlongTokens(ctx, userID.Hex())

		c.Set(userContext, user)
	}
}

func parseAuthorizationHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if headerParts[1] == "" {
		return "", errors.New("access token is empty")
	}

	return headerParts[1], nil
}
