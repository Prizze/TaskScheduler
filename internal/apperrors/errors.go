package apperrors

import "net/http"

type Code string

type Response struct {
	HTTPStatus int
	Code       Code
	Message    string
}

const (
	CodeValidation   Code = "validation_error"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeInternal     Code = "internal_error"
)

var (
	Validation = Response{
		HTTPStatus: http.StatusBadRequest,
		Code:       CodeValidation,
		Message:    "invalid request data",
	}
	Unauthorized = Response{
		HTTPStatus: http.StatusUnauthorized,
		Code:       CodeUnauthorized,
		Message:    "authentication required",
	}
	Forbidden = Response{
		HTTPStatus: http.StatusForbidden,
		Code:       CodeForbidden,
		Message:    "access denied",
	}
	NotFound = Response{
		HTTPStatus: http.StatusNotFound,
		Code:       CodeNotFound,
		Message:    "resource not found",
	}
	Conflict = Response{
		HTTPStatus: http.StatusConflict,
		Code:       CodeConflict,
		Message:    "resource already exists",
	}
	Internal = Response{
		HTTPStatus: http.StatusInternalServerError,
		Code:       CodeInternal,
		Message:    "internal server error",
	}
)
