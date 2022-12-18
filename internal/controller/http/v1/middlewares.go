package v1

import (
	"strings"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kenplix/url-shrtnr/internal/entity"
)

const userContext = "user"

func (h *Handler) userIdentityMiddleware(c *gin.Context) {
	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	accessToken, err := parseAuthorizationHeader(c)
	if err != nil {
		logger.Warn("failed to parse header",
			zap.String("header", "Authorization"),
			zap.Error(err),
		)
		unauthorizedErrorResponse(c)

		return
	}

	claims, err := h.services.JWT.ParseAccessToken(accessToken)
	if err != nil {
		logger.Warn("failed to parse access token",
			zap.String("token", accessToken),
			zap.Error(err),
		)
		unauthorizedErrorResponse(c)

		return
	}

	var g errgroup.Group

	g.Go(func() error {
		e := h.services.JWT.ValidateAccessToken(reqctx, claims)
		if e != nil {
			logger.Warn("failed to validate access token",
				zap.Object("claims", claims),
				zap.Error(e),
			)
			return e
		}

		return nil
	})

	var user entity.User

	g.Go(func() error {
		userID, e := primitive.ObjectIDFromHex(claims.Subject)
		if e != nil {
			logger.Warn("failed to parse userID object",
				zap.String("hex", claims.Subject),
				zap.Error(e),
			)
			return e
		}

		user, e = h.services.Users.GetByID(reqctx, userID)
		if e != nil {
			logger.Warn("failed to get user",
				zap.String("userID", userID.Hex()),
				zap.Error(e),
			)
			return e
		} else if user.SuspendedAt != nil {
			logger.Warn("protected route request from suspended user",
				zap.String("userID", userID.Hex()),
			)
			return &entity.SuspendedUserError{UserID: user.ID.Hex()}
		}

		return nil
	})

	if err = g.Wait(); err != nil {
		var suspUserError *entity.SuspendedUserError
		if errors.As(err, &suspUserError) {
			suspendedErrorResponse(c)
			return
		}

		unauthorizedErrorResponse(c)

		return
	}

	c.Set(userContext, user)
}

func parseAuthorizationHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New(`empty "Authorization" header`)
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New(`invalid "Authorization" header`)
	}

	return headerParts[1], nil
}

func (h *Handler) userActivityMiddleware(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	go h.services.JWT.ProlongTokens(c.Request.Context(), user.ID.Hex())
}
