package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
    return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func NewValidationError(err error) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: "Validation failed",
        Details: err.Error(),
    }
}

func NewNotFoundError(resource string) *AppError {
    return &AppError{
        Code:    http.StatusNotFound,
        Message: fmt.Sprintf("%s not found", resource),
    }
}

func NewInternalError(err error) *AppError {
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: "Internal server error",
        Details: err.Error(),
    }
}

func NewUnauthorizedError() *AppError {
    return &AppError{
        Code:    http.StatusUnauthorized,
        Message: "Unauthorized",
    }
}

func NewForbiddenError() *AppError {
    return &AppError{
        Code:    http.StatusForbidden,
        Message: "Forbidden",
    }
}