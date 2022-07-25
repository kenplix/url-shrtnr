package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/usecase"
)

type usersHandler struct {
	Users usecase.UsersService
}

func NewUsersHandler(users usecase.UsersService) *usersHandler {
	return &usersHandler{
		Users: users,
	}
}

func (h *usersHandler) initRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")

	users.POST("/sign-up", h.userSignUp)
	users.POST("/sign-in", h.userSignIn)
}

type userSignUpInput struct {
	FirstName string `json:"firstName" binding:"required,min=2,max=32"`
	LastName  string `json:"lastName" binding:"required,min=2,max=32"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=64"`
}

func (h *usersHandler) userSignUp(c *gin.Context) {
	var input userSignUpInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Message: errInvalidInputBody.Error(),
		})

		return
	}

	err := h.Users.SignUp(c.Request.Context(), usecase.UserSignUpInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
				Message: entity.ErrUserAlreadyExists.Error(),
			})

			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.Status(http.StatusCreated)
}

type userSignInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *usersHandler) userSignIn(c *gin.Context) {
	var input userSignInInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Message: errInvalidInputBody.Error(),
		})

		return
	}

	tokens, err := h.Users.SignIn(c.Request.Context(), usecase.UserSignInInput{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrIncorrectEmailOrPassword) {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
				Message: entity.ErrIncorrectEmailOrPassword.Error(),
			})

			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
