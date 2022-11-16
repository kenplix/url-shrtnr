package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Kenplix/url-shrtnr/internal/entity"
	"github.com/Kenplix/url-shrtnr/internal/service"
	"github.com/Kenplix/url-shrtnr/pkg/auth"
)

const userContext = "userID"

type usersHandler struct {
	usersServ  service.UsersService
	tokensServ auth.TokensService
}

func NewUsersHandler(usersServ service.UsersService, tokensServ auth.TokensService) *usersHandler {
	return &usersHandler{
		usersServ:  usersServ,
		tokensServ: tokensServ,
	}
}

func (h *usersHandler) initRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/sign-up", h.userSignUp)
		users.POST("/sign-in", h.userSignIn)
		users.POST("/refresh-tokens", h.userRefreshTokens)

		authenticated := users.Group("/", h.userIdentity)
		authenticated.GET("/protected", h.userProtectedRoute)
	}
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

	err := h.usersServ.SignUp(c.Request.Context(), service.UserSignUpInput{
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

func (h *usersHandler) userSignIn(c *gin.Context) {
	var input userSignInInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Message: errInvalidInputBody.Error(),
		})

		return
	}

	tokens, err := h.usersServ.SignIn(c.Request.Context(), service.UserSignInInput{
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

	c.JSON(http.StatusOK, tokens)
}

type userRefreshTokensInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (h *usersHandler) userRefreshTokens(c *gin.Context) {
	var input userRefreshTokensInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Message: errInvalidInputBody.Error(),
		})

		return
	}

	tokens, err := h.usersServ.RefreshTokens(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *usersHandler) userIdentity(c *gin.Context) {
	userID, err := h.parseAuthorizationHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
			Message: err.Error(),
		})

		return
	}

	c.Set(userContext, userID)
}

func (h *usersHandler) parseAuthorizationHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if headerParts[1] == "" {
		return "", errors.New("token is empty")
	}

	return h.tokensServ.ParseAccessToken(headerParts[1])
}

func getUserID(c *gin.Context) (primitive.ObjectID, error) {
	return getIDByContext(c, userContext)
}

func getIDByContext(c *gin.Context, context string) (primitive.ObjectID, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("%q context not found", context)
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return primitive.ObjectID{}, fmt.Errorf("%q context is of invalid type", context)
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return id, nil
}

func (h *usersHandler) userProtectedRoute(c *gin.Context) {
	id, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
		})

		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"userID": id.Hex(),
	})
}
