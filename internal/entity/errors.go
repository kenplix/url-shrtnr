package entity

import "errors"

var (
	ErrUserNotFound             = errors.New("user doesn't exists")
	ErrUserAlreadyExists        = errors.New("user with such email already exists")
	ErrIncorrectEmailOrPassword = errors.New("incorrect email or password")
)
