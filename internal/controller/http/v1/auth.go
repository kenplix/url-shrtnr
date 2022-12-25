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
	Username string `json:"username" binding:"required,username" example:"kenplix"`
	Email    string `json:"email" binding:"required,email" example:"tolstoi.job@gmail.com"`
	Password string `json:"password" binding:"required,password" example:"1wE$Rty2"`
}

// @Summary		Registers users accounts
// @Tags			auth
// @Description	Registers users accounts
// @Accept			json
// @Produce		json
// @Param			schema	body	userSignUpSchema	true	"JSON schema for user account registration"
// @Success		201		"User account was successfully registered"
// @Failure		400		{object}	errResponse{errors=[]entity.CoreError}			"Invalid JSON or wrong type of JSON values"
// @Failure		422		{object}	errResponse{errors=[]entity.ValidationError}	"Validation failed through invalid fields"
// @Failure		500		{object}	errResponse{errors=[]entity.CoreError}			"Internal server error"
// @Router			/auth/sign-up [post]
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
	Login    string `json:"login" binding:"required,login" example:"kenplix or tolstoi.job@gmail.com"`
	Password string `json:"password" binding:"required,password" example:"1wE$Rty2"`
}

// @Summary		Logins users accounts
// @Tags			auth
// @Description	Logins users accounts
// @Accept			json
// @Produce		json
// @Param			schema	body		userSignInSchema								true	"JSON schema for user login"
// @Success		200		{object}	entity.Tokens									"User was successfully logged in"
// @Failure		400		{object}	errResponse{errors=[]entity.CoreError}			"Invalid JSON or wrong type of JSON values"
// @Failure		403		{object}	errResponse{errors=[]entity.CoreError}			"Your account has been suspended"
// @Failure		422		{object}	errResponse{errors=[]entity.ValidationError}	"Validation failed through invalid fields"
// @Failure		500		{object}	errResponse{errors=[]entity.CoreError}			"Internal server error"
// @Router			/auth/sign-in [post]
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

// @Summary		Logout users from the server
// @Security		JWT-RS256
// @Tags			auth
// @Description	Logout users from the server
// @Produce		json
// @Success		200	"User was successfully signed out"
// @Failure		401	{object}	errResponse{errors=[]entity.CoreError}	"Access is denied due to invalid credentials"
// @Failure		403	{object}	errResponse{errors=[]entity.CoreError}	"Your account has been suspended"
// @Failure		500	{object}	errResponse{errors=[]entity.CoreError}	"Internal server error"
// @Router			/auth/sign-out [post]
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
	RefreshToken string `json:"refreshToken" binding:"required,jwt" example:"header.payload.signature"`
}

// @Summary		Refresh users tokens
// @Tags			auth
// @Description	Refresh users tokens
// @Accept			json
// @Produce		json
// @Param			schema	body		userRefreshTokensSchema							true	"JSON schema for tokens refresh"
// @Success		200		{object}	entity.Tokens									"User tokens was successfully refreshed"
// @Failure		400		{object}	errResponse{errors=[]entity.CoreError}			"Invalid JSON or wrong type of JSON values"
// @Failure		422		{object}	errResponse{errors=[]entity.ValidationError}	"Validation failed through invalid fields"
// @Failure		500		{object}	errResponse{errors=[]entity.CoreError}			"Internal server error"
// @Router			/auth/refresh-tokens [post]
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
