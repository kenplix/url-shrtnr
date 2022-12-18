package v1

import (
	"net/http"
	"strings"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/kenplix/url-shrtnr/internal/service"

	"github.com/kenplix/url-shrtnr/internal/entity"
)

func (h *Handler) initAuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")

	auth.POST("/sign-up", h.signUp)
	auth.POST("/sign-in", h.signIn)
	auth.POST("/sign-out", h.userIdentityMiddleware, h.signOut)
	auth.POST("/refresh-tokens", h.refreshTokens)
}

type userSignUpSchema struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

func (h *Handler) signUp(c *gin.Context) {
	var schema userSignUpSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	err := h.services.Auth.SignUp(reqctx, service.UserSignUpSchema{
		Username: schema.Username,
		Email:    strings.ToLower(schema.Email),
		Password: schema.Password,
	})
	if err != nil {
		var validationError *entity.ValidationError
		if errors.As(err, &validationError) {
			logger.Warn("failed to sign up", zap.Error(err))
			errorResponse(c, http.StatusUnprocessableEntity, validationError)

			return
		}

		logger.Error("failed to sign up", zap.Error(err))
		internalErrorResponse(c)

		return
	}

	c.Status(http.StatusCreated)
}

type userSignInSchema struct {
	Login    string `json:"login" binding:"required,login"`
	Password string `json:"password" binding:"required,password"`
}

func (h *Handler) signIn(c *gin.Context) {
	var schema userSignInSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	tokens, err := h.services.Auth.SignIn(reqctx, service.UserSignInSchema{
		Login:    schema.Login,
		Password: schema.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrIncorrectCredentials) {
			logger.Warn("failed to sign in", zap.Error(err))
			errorResponse(c, http.StatusUnprocessableEntity, &entity.CoreError{
				Code:    errorcode.IncorrectCredentials,
				Message: entity.ErrIncorrectCredentials.Error(),
			})

			return
		}

		var suspUserError *entity.SuspendedUserError
		if errors.As(err, &suspUserError) {
			logger.Debug("suspended user tries to sign in",
				zap.String("userID", suspUserError.UserID),
			)
			suspendedErrorResponse(c)

			return
		}

		logger.Error("failed to sign in", zap.Error(err))
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) signOut(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	err := h.services.Auth.SignOut(reqctx, user.ID)
	if err != nil {
		logger.Error("failed to sign out",
			zap.String("userID", user.ID.Hex()),
			zap.Error(err),
		)
		internalErrorResponse(c)

		return
	}

	c.Status(http.StatusOK)
}

type userRefreshTokensSchema struct {
	RefreshToken string `json:"refreshToken" binding:"required,jwt"`
}

func (h *Handler) refreshTokens(c *gin.Context) {
	var schema userRefreshTokensSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	claims, err := h.services.JWT.ParseRefreshToken(schema.RefreshToken)
	if err != nil {
		logger.Warn("failed to parse refresh token",
			zap.String("token", schema.RefreshToken),
			zap.Error(err),
		)
		errorResponse(c, http.StatusUnprocessableEntity, &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.InvalidField,
				Message: "refresh token is invalid, expired or revoked",
			},
			Field: "refreshToken",
		})

		return
	}

	err = h.services.JWT.ValidateRefreshToken(reqctx, claims)
	if err != nil {
		logger.Warn("failed to validate refresh token",
			zap.String("token", schema.RefreshToken),
			zap.Error(err),
		)
		errorResponse(c, http.StatusUnprocessableEntity, &entity.ValidationError{
			CoreError: entity.CoreError{
				Code:    errorcode.InvalidField,
				Message: "refresh token is invalid, expired or revoked",
			},
			Field: "refreshToken",
		})

		return
	}

	tokens, err := h.services.JWT.CreateTokens(reqctx, claims.Subject)
	if err != nil {
		logger.Error("failed to create tokens pair",
			zap.String("userID", claims.Subject),
			zap.Error(err),
		)
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}
