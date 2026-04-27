package errs

import (
	"errors"
	"net/http"
)

type Code int

const (
	CodeOK              Code = 0
	CodeBadRequest      Code = 4000
	CodeUnauthorized    Code = 4001
	CodeForbidden       Code = 4003
	CodeNotFound        Code = 4004
	CodeConflict        Code = 4009
	CodeRateLimited     Code = 4029
	CodeInternal        Code = 5000
	CodeInsufficientSKU Code = 6001
	CodeDuplicate       Code = 6002
)

type AppError struct {
	Code    Code
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(code Code, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code Code, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause}
}

func From(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Wrap(CodeInternal, "internal_error", err)
}

func ToHTTP(err error) (int, int, string) {
	appErr := From(err)
	switch appErr.Code {
	case CodeBadRequest:
		return http.StatusBadRequest, int(appErr.Code), appErr.Message
	case CodeUnauthorized:
		return http.StatusUnauthorized, int(appErr.Code), appErr.Message
	case CodeForbidden:
		return http.StatusForbidden, int(appErr.Code), appErr.Message
	case CodeNotFound:
		return http.StatusNotFound, int(appErr.Code), appErr.Message
	case CodeConflict, CodeDuplicate:
		return http.StatusConflict, int(appErr.Code), appErr.Message
	case CodeRateLimited:
		return http.StatusTooManyRequests, int(appErr.Code), appErr.Message
	case CodeInsufficientSKU:
		return http.StatusBadRequest, int(appErr.Code), appErr.Message
	default:
		return http.StatusInternalServerError, int(CodeInternal), "internal_error"
	}
}
