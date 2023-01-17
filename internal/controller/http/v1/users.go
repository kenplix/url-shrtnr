package v1

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"

	"go.uber.org/zap"

	"github.com/kenplix/url-shrtnr/internal/service"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"github.com/kenplix/url-shrtnr/internal/entity"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(router *gin.RouterGroup) {
	users := router.Group(
		"/users",
		h.userIdentityMiddleware,
		h.userActivityMiddleware,
	)

	users.GET("/me", h.me)
	users.PATCH("/change-email", h.changeEmail)
	users.PATCH("/change-password", h.changePassword)
}

// me handler returns users personal information
//
//	@Summary		Returns users personal information
//	@Security		JWT-RS256
//	@Tags			user
//	@Description	Returns users personal information
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entity.User								"User personal information"
//	@Failure		401	{object}	errResponse{errors=[]entity.CoreError}	"Access is denied due to invalid credentials"
//	@Failure		403	{object}	errResponse{errors=[]entity.CoreError}	"Your account has been suspended"
//	@Failure		500	{object}	errResponse{errors=[]entity.CoreError}	"Internal server error"
//	@Router			/users/me [get]
func (h *Handler) me(c *gin.Context) {
	user := c.MustGet(userContext).(entity.User)

	c.JSON(http.StatusOK, user)
}

type userChangeEmailSchema struct {
	NewEmail string `json:"newEmail" binding:"required,email" example:"example@gmail.com"`
}

// changeEmail handler changes users emails
//
//	@Summary		Changes users emails
//	@Security		JWT-RS256
//	@Tags			user
//	@Description	Changes users emails
//	@Accept			json
//	@Produce		json
//	@Param			schema	body	userChangeEmailSchema	true	"JSON schema for user email changing"
//	@Success		200		"User email was successfully changed"
//	@Failure		400		{object}	errResponse{errors=[]entity.CoreError}			"Invalid JSON or wrong type of JSON values"
//	@Failure		401		{object}	errResponse{errors=[]entity.CoreError}			"Access is denied due to invalid credentials"
//	@Failure		403		{object}	errResponse{errors=[]entity.CoreError}			"Your account has been suspended"
//	@Failure		422		{object}	errResponse{errors=[]entity.ValidationError}	"Validation failed through invalid fields"
//	@Failure		500		{object}	errResponse{errors=[]entity.CoreError}			"Internal server error"
//	@Router			/users/change-email [patch]
func (h *Handler) changeEmail(c *gin.Context) {
	var schema userChangeEmailSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	user := c.MustGet(userContext).(entity.User)

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	err := h.services.Users.ChangeEmail(reqctx, service.ChangeEmailSchema{
		UserID:   user.ID,
		NewEmail: schema.NewEmail,
	})
	if err != nil {
		logger.Error("failed to change email",
			zap.String("userID", user.ID.Hex()),
			zap.Error(err),
		)
		internalErrorResponse(c)

		return
	}

	c.Status(http.StatusOK)
}

type userChangePasswordSchema struct {
	CurrentPassword      string `json:"currentPassword" binding:"required,password" example:"1wE$Rty2"`
	NewPassword          string `json:"newPassword" binding:"required,password,eqfield=PasswordConfirmation" example:"2ytR$Ew1"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,password" example:"2ytR$Ew1"`
}

// changePassword handler changes users passwords
//
//	@Summary		Changes users passwords
//	@Security		JWT-RS256
//	@Tags			user
//	@Description	Changes users passwords
//	@Accept			json
//	@Produce		json
//	@Param			schema	body	userChangePasswordSchema	true	"JSON schema for user password changing"
//	@Success		200		"User password was successfully changed"
//	@Failure		400		{object}	errResponse{errors=[]entity.CoreError}			"Invalid JSON or wrong type of JSON values"
//	@Failure		401		{object}	errResponse{errors=[]entity.CoreError}			"Access is denied due to invalid credentials"
//	@Failure		403		{object}	errResponse{errors=[]entity.CoreError}			"Your account has been suspended"
//	@Failure		422		{object}	errResponse{errors=[]entity.ValidationError}	"Validation failed through invalid fields"
//	@Failure		500		{object}	errResponse{errors=[]entity.CoreError}			"Internal server error"
//	@Router			/users/change-password [patch]
func (h *Handler) changePassword(c *gin.Context) {
	var schema userChangePasswordSchema
	if err := c.ShouldBindJSON(&schema); err != nil {
		bindingErrorResponse(c, err)
		return
	}

	user := c.MustGet(userContext).(entity.User)

	reqctx := c.Request.Context()
	logger := log.LoggerFromContext(reqctx)

	err := h.services.Users.ChangePassword(reqctx, service.ChangePasswordSchema{
		UserID:          user.ID,
		CurrentPassword: schema.CurrentPassword,
		NewPassword:     schema.NewPassword,
	})
	if err != nil {
		if errors.Is(err, entity.ErrIncorrectCredentials) {
			logger.Warn("failed to change password",
				zap.String("userID", user.ID.Hex()),
				zap.Error(err),
			)
			errorResponse(c, http.StatusUnprocessableEntity, &entity.CoreError{
				Code:    errorcode.IncorrectCredentials,
				Message: entity.ErrIncorrectCredentials.Error(),
			})

			return
		}

		logger.Error("failed to change password",
			zap.String("userID", user.ID.Hex()),
			zap.Error(err),
		)
		internalErrorResponse(c)

		return
	}

	c.Status(http.StatusOK)
}
