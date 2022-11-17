package entity

import (
	"errors"
	"fmt"
	"github.com/Kenplix/url-shrtnr/internal/entity/errorcode"
)

var (
	ErrUserNotFound         = errors.New("user doesn't exists")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrIncorrectCredentials = errors.New("incorrect credentials")
)

type CoreError struct {
	Code    errorcode.ErrorCode `json:"code"`
	Message string              `json:"message"`
}

func (e *CoreError) ErrorCode() errorcode.ErrorCode { return e.Code }
func (e *CoreError) ErrorMessage() string           { return e.Message }
func (e *CoreError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.Code, e.Message)
}

type ValidationError struct {
	CoreError
	Field string `json:"field"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s [%s]: %s", e.Field, e.Code, e.Message)
}
