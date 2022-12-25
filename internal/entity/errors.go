package entity

import (
	"errors"
	"fmt"

	"github.com/kenplix/url-shrtnr/internal/entity/errorcode"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

type SuspendedUserError struct {
	UserID string
}

func (e *SuspendedUserError) Error() string {
	return fmt.Sprintf("user[id:%q] suspended", e.UserID)
}

// CoreError is a basic representation of API call error
//
//	@Description	Basic representation of API call error
type CoreError struct {
	// Code is CAPS_CASE constant error code you can programmatically consume to make resolution decisions from
	Code errorcode.ErrorCode `json:"code" example:"ERROR_CODE"`
	// Message indicate a (usually) human-readable description of the error
	Message string `json:"message" example:"error cause description"`
}

func (e *CoreError) ErrorCode() errorcode.ErrorCode { return e.Code }
func (e *CoreError) ErrorMessage() string           { return e.Message }
func (e *CoreError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

// ValidationError is a standardized representation of a validation errors
//
//	@Description	Standardized representation of a validation errors
type ValidationError struct {
	CoreError
	// Field with which validation error related
	Field string `json:"field" example:"invalid field"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s [%s]: %s", e.Field, e.Code, e.Message)
}
