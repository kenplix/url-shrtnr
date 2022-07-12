package v1

import "errors"

var errInvalidInputBody = errors.New("invalid input body")

type ErrorResponse struct {
	Message string `json:"message"`
}
