package v1

import (
	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
)

type AuthHandler struct {
	authServ   service.AuthService
	tokensServ auth.TokensService
}

func NewAuthHandler(authServ service.AuthService, tokensServ auth.TokensService) (*AuthHandler, error) {
	if authServ == nil {
		return nil, errors.New("auth service not provided")
	}

	if tokensServ == nil {
		return nil, errors.New("tokens service not provided")
	}

	h := &AuthHandler{
		authServ:   authServ,
		tokensServ: tokensServ,
	}

	return h, nil
}

func (h *AuthHandler) init(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/sign-up", h.signUp)
		authGroup.POST("/sign-in", h.signIn)
		authGroup.POST("/refresh-tokens", h.refreshTokens)
	}
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
		errorResponse(c, http.StatusInternalServerError, &entity.CoreError{
			Code:    errorcode.InternalError,
			Message: http.StatusText(http.StatusInternalServerError),
		})
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

		log.Printf("error: failed to sign in: %s", err)
		errorResponse(c, http.StatusInternalServerError, &entity.CoreError{
			Code:    errorcode.InternalError,
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
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

	userID, err := h.tokensServ.ParseRefreshToken(schema.RefreshToken)
	if err != nil {
		log.Printf("warning: failed to parse %q refresh token: %s", schema.RefreshToken, err)
		errorResponse(c, http.StatusBadRequest, &entity.CoreError{
			Code:    errorcode.ParsingError,
			Message: "problems parsing JWT",
		})
		return
	}

	tokens, err := h.tokensServ.CreateTokens(userID)
	if err != nil {
		log.Printf("error: failed to create user[id:%s] tokens: %s", userID, err)
		errorResponse(c, http.StatusInternalServerError, &entity.CoreError{
			Code:    errorcode.InternalError,
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
