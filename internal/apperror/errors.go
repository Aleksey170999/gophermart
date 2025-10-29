package apperror

import (
	"net/http"
)

type AppError struct {
	Err        error
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrInvalidOrderNumber = &AppError{
		Message:    "invalid order number format",
		HTTPStatus: http.StatusUnprocessableEntity,
	}

	ErrOrderByUserExists = &AppError{
		Message:    "order already uploaded by this user",
		HTTPStatus: http.StatusOK,
	}

	ErrOrderByOtherExists = &AppError{
		Message:    "order already uploaded by another user",
		HTTPStatus: http.StatusConflict,
	}

	ErrUnauthorized = &AppError{
		Message:    "unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrInternal = &AppError{
		Message:    "internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
)

func New(message string, status int) *AppError {
	return &AppError{
		Message:    message,
		HTTPStatus: status,
	}
}

func Wrap(err error, message string, status int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		HTTPStatus: status,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
