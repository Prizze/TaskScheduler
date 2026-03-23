package domain

import "errors"

var (
	ErrServerError     = errors.New("internal server error")
	ErrNoTag           = errors.New("tag id is not exist")
	ErrInvalidStatus   = errors.New("invalid task status")
	ErrInvalidPriority = errors.New("invalid task priority")
	ErrValidation      = errors.New("validation error")
)