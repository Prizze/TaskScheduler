package domain

import (
	"errors"
)

var (
	ErrServerError        = errors.New("internal error")
	ErrValidation         = errors.New("validation error")
	ErrEmailIsTaken       = errors.New("email is taken")
	ErrUserNotFound       = errors.New("user is not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
