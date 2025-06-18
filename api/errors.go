package api

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type APIError struct {
	Error   error
	Code    string
	Message string
}

func NewError(code string, message string) APIError {
	return APIError{
		Error:   errors.New(message),
		Code:    code,
		Message: message,
	}
}

var (
	ErrHTTPBadRequest             = NewError("ERR_HTTP_BADREQUEST", "bad request")
	ErrHTTPInternal               = NewError("ERR_HTTP_INTERNAL", "something went wrong :(")
	ErrHTTPNotFound               = NewError("ERR_HTTP_NOTFOUND", "resource not found")
	ErrHTTPForbidden              = NewError("ERR_HTTP_FORBIDDEN", "you cannot use this route")
	ErrJSONValidation             = NewError("ERR_JSON_VALIDATION", "validation failed")
	ErrUnknownEmail               = NewError("ERR_AUTH_EMAIL_UNKNOWN", "unknown email")
	ErrWrongPassword              = NewError("ERR_AUTH_WRONG_PASSWORD", "wrong password")
	ErrUnexpectedJWTSigningMethod = NewError("ERR_AUTH_UNEXPECTED_JWT_SIGN_METHOD", "unexpected signing method")
	ErrEmailAlreadyExists         = NewError("ERR_AUTH_EMAIL_CONFLICT", "email conflict")
	ErrJWTInvalid                 = NewError("ERR_JWT_INVALID", "invalid jwt")
	ErrLoginInvalid               = NewError("ERR_LOGIN_INVALIDCREDS", "invalid credentials")
	ErrRegistrationFailed         = NewError("ERR_REGISTER_FAILED", "registration failed")
	ErrInvalidUUID                = NewError("ERR_INVALID_UUID", "invalid uuid")
	ErrUnknownTable               = NewError("ERR_DINING_UNKNOWNTABLE", "table does not exist")
	ErrTableNotAvaliable          = NewError("ERR_DINING_TABLE_UNAVAILABLE", "table is not available")
	ErrUnknownOrder               = NewError("ERR_ORDER_UNKNOWN", "order does not exist")
	ErrOrderNotOngoing            = NewError("ERR_ORDER_NOTONGOING", "order is not ongoing")
	ErrUnknownOrderItem           = NewError("ERR_ORDER_UNKNOWNITEM", "order item does not exist")
	ErrItemNameConflict           = NewError("ERR_MENU_ITEMNAME_CONFLICT", "menu item with the same name already exists")
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func NewErrorResponse(status int, code string, msg string, errs any) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Code:    code,
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
			field := strings.ToLower(e.Field())
			tag := e.Tag()

			res = append(res, NewValidationError(tag, field))
		}
	}
	return res
}
