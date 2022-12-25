package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kenplix/url-shrtnr/internal/controller/http/ginctx"

	"github.com/kenplix/url-shrtnr/pkg/log"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"github.com/kenplix/url-shrtnr/internal/entity"
	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"
)

type apiError interface {
	ErrorCode() errorcode.ErrorCode
	ErrorMessage() string
}

// errResponse is a standardized representation of an errors that may occur in API calls
//
//	@Description	Standardized representation of an errors that may occur in API calls
type errResponse struct {
	// Array of errors that occurred performing API call
	// Usually have one error inside, except validation errors where we have error for each invalid field
	Errors []apiError `json:"errors"`
}

func errorResponse(c *gin.Context, code int, apiErrors ...apiError) {
	if code < 400 || len(apiErrors) == 0 {
		logger := log.LoggerFromContext(c.Request.Context())
		logger.Warn("calling errorResponse function without errors",
			zap.Int("response-code", code),
			zap.Int("errors-count", len(apiErrors)),
		)

		return
	}

	c.AbortWithStatusJSON(code, errResponse{Errors: apiErrors})
}

func suspendedErrorResponse(c *gin.Context) {
	errorResponse(c, http.StatusForbidden, newSuspendedError())
}

func newSuspendedError() *entity.CoreError {
	return &entity.CoreError{
		Code:    errorcode.CurrentUserSuspended,
		Message: "your account has been suspended",
	}
}

func unauthorizedErrorResponse(c *gin.Context) {
	errorResponse(c, http.StatusUnauthorized, newUnauthorizedError())
}

func newUnauthorizedError() *entity.CoreError {
	return &entity.CoreError{
		Code:    errorcode.UnauthorizedAccess,
		Message: "access is denied due to invalid credentials",
	}
}

func internalErrorResponse(c *gin.Context) {
	errorResponse(c, http.StatusInternalServerError, newInternalError())
}

func newInternalError() *entity.CoreError {
	return &entity.CoreError{
		Code:    errorcode.InternalError,
		Message: strings.ToLower(http.StatusText(http.StatusInternalServerError)),
	}
}

func bindingErrorResponse(c *gin.Context, err error) {
	logger := log.LoggerFromContext(c.Request.Context())

	if err == nil {
		logger.Warn("calling bindingErrorResponse function without error")
		return
	}

	logger.Warn("request binding error", zap.Error(err))

	switch typedError := err.(type) {
	case validator.ValidationErrors:
		validationErrors := make([]apiError, 0, len(typedError))
		for _, fieldError := range typedError {
			validationErrors = append(validationErrors, parseFieldError(c, fieldError))
		}

		errorResponse(c, http.StatusUnprocessableEntity, validationErrors...)
	case *json.UnmarshalTypeError:
		errorResponse(c, http.StatusBadRequest, &entity.CoreError{
			Code:    errorcode.InvalidSchema,
			Message: parseUnmarshalTypeError(*typedError),
		})
	default:
		errorResponse(c, http.StatusBadRequest, &entity.CoreError{
			Code:    errorcode.ParsingError,
			Message: "problems parsing JSON",
		})
	}
}

func parseFieldError(c *gin.Context, err validator.FieldError) *entity.ValidationError {
	code := errorcode.InvalidField

	if err.Tag() == "required" {
		code = errorcode.MissingField
	}

	translator := c.MustGet(ginctx.TranslatorContext).(ut.Translator)

	return &entity.ValidationError{
		CoreError: entity.CoreError{
			Code:    code,
			Message: err.Translate(translator),
		},
		Field: err.Field(),
	}
}

func parseUnmarshalTypeError(err json.UnmarshalTypeError) string {
	if err.Field != "" {
		return fmt.Sprintf("%s field must be a %s type", err.Field, err.Type.String())
	}

	return "body should be a JSON object"
}
