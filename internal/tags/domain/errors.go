package domain

import "errors"

var (
	ErrServerError      = errors.New("internal server error")
	ErrValidation       = errors.New("validation error")
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagAlreadyExists = errors.New("tag already exists")
)
