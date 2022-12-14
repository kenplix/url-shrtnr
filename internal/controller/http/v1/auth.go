package v1

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/service"

	"github.com/Kenplix/url-shrtnr/internal/entity"
)

type AuthHandler struct {
	authServ  service.AuthService
	usersServ service.UsersService
	jwtServ   service.JWTService
}

func NewAuthHandler(
	authServ service.AuthService,
	usersServ service.UsersService,
	jwtServ service.JWTService,
) (*AuthHandler, error) {
	if authServ == nil {
		return nil, errors.New("auth service not provided")
	}

	if usersServ == nil {
		return nil, errors.New("users service not provided")
	}

	if jwtServ == nil {
		return nil, errors.New("jwt service not provided")
	}

	h := &AuthHandler{
		authServ:  authServ,
		usersServ: usersServ,
		jwtServ:   jwtServ,
	}

	return h, nil
}

func (h *AuthHandler) init(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")

	authGroup.POST("/sign-up", h.signUp)
	authGroup.POST("/sign-in", h.signIn)
	authGroup.POST("/sign-out", userIdentityMiddleware(h.usersServ, h.jwtServ), h.signOut)
	authGroup.POST("/refresh-tokens", h.refreshTokens)
}

type userSignUpSchema struct {
	Username string `json:"username" binding:"required,username"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,password"`
}

func (h *AuthHandler) signUp(c *gin.Context) {
	var schema userSignUpSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	err := h.authServ.SignUp(c.Request.Context(), service.UserSignUpSchema{
		Username: schema.Username,
		Email:    strings.ToLower(schema.Email),
		Password: schema.Password,
	})
	if err != nil {
		var validationError *entity.ValidationError
		if errors.As(err, &validationError) {
			zap.L().Warn("failed to sign up", zap.Error(err))
			errorResponse(c, http.StatusUnprocessableEntity, validationError)

			return
		}

		zap.L().Error("failed to sign up", zap.Error(err))
		internalErrorResponse(c)

		return
	}

	c.Status(http.StatusCreated)
}

type userSignInSchema struct {
	Login    string `json:"login" binding:"required,login"`
	Password string `json:"password" binding:"required,password"`
}

func (h *AuthHandler) signIn(c *gin.Context) {
	var schema userSignInSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	tokens, err := h.authServ.SignIn(c.Request.Context(), service.UserSignInSchema{
		Login:    schema.Login,
		Password: schema.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrIncorrectCredentials) {
			zap.L().Warn("failed to sign in", zap.Error(err))
			errorResponse(c, http.StatusUnprocessableEntity, &entity.CoreError{
				Code:    errorcode.IncorrectCredentials,
				Message: entity.ErrIncorrectCredentials.Error(),
			})

			return
		}

		var suspUserError *entity.SuspendedUserError
		if errors.As(err, &suspUserError) {
			zap.L().Debug("suspended user tries to sign in",
				zap.String("userID", suspUserError.UserID),
			)
			suspendedErrorResponse(c)

			return
		}

		zap.L().Error("failed to sign in", zap.Error(err))
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) signOut(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	err := h.authServ.SignOut(c.Request.Context(), user.ID)
	if err != nil {
		zap.L().Error("failed to sign out",
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

func (h *AuthHandler) refreshTokens(c *gin.Context) {
	var schema userRefreshTokensSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	claims, err := h.jwtServ.ParseRefreshToken(schema.RefreshToken)
	if err != nil {
		zap.L().Warn("failed to parse refresh token",
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

	err = h.jwtServ.ValidateRefreshToken(c.Request.Context(), claims)
	if err != nil {
		zap.L().Warn("failed to validate refresh token",
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

	tokens, err := h.jwtServ.CreateTokens(c.Request.Context(), claims.Subject)
	if err != nil {
		zap.L().Error("failed to create tokens pair",
			zap.String("userID", claims.Subject),
			zap.Error(err),
		)
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}
