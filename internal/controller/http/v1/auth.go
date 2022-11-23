package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/service"

	"github.com/Kenplix/url-shrtnr/internal/entity"
)

type AuthHandler struct {
	authServ service.AuthService
	jwtServ  service.JWTService
}

func NewAuthHandler(authServ service.AuthService, jwtServ service.JWTService) (*AuthHandler, error) {
	if authServ == nil {
		return nil, errors.New("auth service not provided")
	}

	if jwtServ == nil {
		return nil, errors.New("jwt service not provided")
	}

	h := &AuthHandler{
		authServ: authServ,
		jwtServ:  jwtServ,
	}

	return h, nil
}

func (h *AuthHandler) init(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")

	authGroup.POST("/sign-up", h.signUp)
	authGroup.POST("/sign-in", h.signIn)
	authGroup.POST("/sign-out", userIdentityMiddleware(h.jwtServ), h.signOut)
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
			log.Printf("warning: failed to sign up: %s", err)
			errorResponse(c, http.StatusUnprocessableEntity, validationError)

			return
		}

		log.Printf("error: failed to sign up: %s", err)
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
			log.Printf("warning: failed to sign in: %s", err)
			errorResponse(c, http.StatusUnprocessableEntity, &entity.CoreError{
				Code:    errorcode.IncorrectCredentials,
				Message: entity.ErrIncorrectCredentials.Error(),
			})

			return
		}

		var suspUserError *entity.SuspendedUserError
		if errors.As(err, &suspUserError) {
			log.Printf("debug: suspended user[id:%q] tries to sign in", suspUserError.UserID)
			errorResponse(c, http.StatusForbidden, &entity.CoreError{
				Code:    errorcode.CurrentUserSuspended,
				Message: "your account has been suspended",
			})

			return
		}

		log.Printf("error: failed to sign in: %s", err)
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) signOut(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		log.Printf("error: failed to get userID object: %s", err)
		internalErrorResponse(c)

		return
	}

	err = h.authServ.SignOut(c.Request.Context(), userID)
	if err != nil {
		log.Printf("error: user[id:%q]: failed to sign out: %s", userID, err)
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
		log.Printf("warning: failed to parse %q refresh token: %s", schema.RefreshToken, err)
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
		log.Printf("warning: failed to validate %q refresh token: %s", schema.RefreshToken, err)
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
		log.Printf("error: user[id:%q]: failed to create tokens: %s", claims.Subject, err)
		internalErrorResponse(c)

		return
	}

	c.JSON(http.StatusOK, tokens)
}
