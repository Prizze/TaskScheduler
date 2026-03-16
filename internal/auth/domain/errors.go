package domain

import (
	"errors"
)

var (
	ErrEmailIsTaken = errors.New("email is taken")
)
