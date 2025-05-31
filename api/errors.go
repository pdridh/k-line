package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"detail"`
	Errors  any    `json:"errors,omitempty"`
}

func NewErrorResponse(status int, msg string, errs any) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: msg,
		Errors:  errs,
	}
}

type ValidationError struct {
	Code  string `json:"code"`
	Field string `json:"field"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s on %s", e.Code, e.Field)
}

func NewValidationError(tag string, field string) ValidationError {
	return ValidationError{
		Code:  tag, // TODO change this to custom error codes like ErrMissingField based on the tag
		Field: field,
	}
}

func FormatValidationErrors(err error) []ValidationError {
	var res []ValidationError

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			tag := e.Tag()

			res = append(res, NewValidationError(tag, field))
		}
	}
	return res
}
