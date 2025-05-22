package pkg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

// Application error codes.
//
// NOTE: These are meant to be generic and they map well to HTTP error codes.
// Different applications can have very different error code requirements so
// these should be expanded as needed (or introduce subcodes).
const (
	ECONFLICT            = "conflict"
	EINTERNAL            = "internal"
	EINVALID             = "invalid"
	ENOTFOUND            = "not_found"
	ENOTIMPLEMENTED      = "not_implemented"
	EUNAUTHORIZED        = "unauthorized"
	EUNPROCESSABLEENTITY = "unprocessable_entity"
	ETOOMANYREQUESTS     = "too_many_requests"
	EFORBIDDEN           = "forbidden"
	ECANCELED            = "canceled"
)

var ErrRecordNotFound = Errorf(ENOTFOUND, "Record not found.")

// Lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	ECONFLICT:            http.StatusConflict,
	EINVALID:             http.StatusBadRequest,
	ENOTFOUND:            http.StatusNotFound,
	ENOTIMPLEMENTED:      http.StatusNotImplemented,
	EUNAUTHORIZED:        http.StatusUnauthorized,
	EUNPROCESSABLEENTITY: http.StatusUnprocessableEntity,
	EINTERNAL:            http.StatusInternalServerError,
	ETOOMANYREQUESTS:     http.StatusTooManyRequests,
	EFORBIDDEN:           http.StatusForbidden,
	ECANCELED:            499,
}

// WriteError represents an application-specific error. Application errors can be
// unwrapped by the caller to extract out the code & message.
//
// Any non-application error (such as a disk error) should be reported as an
// EINTERNAL error and the human user should only see "Internal error" as the
// message. These low-level internal error details should only be logged and
// reported to the operator of the application (not the end user).
type Error struct {
	// Optional wrapped error.
	Err error

	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

// WriteError implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("app error: code=%s message=%s", e.Code, e.Message)
}

// Is implements the error comparison based on the error code.
func (e *Error) Is(err error) bool {
	var arhErr *Error
	if errors.As(err, &arhErr) {
		return arhErr.Code == e.Code
	}
	return false
}

// LogValue implements slog.LogValuer.
// It returns a group containing the fields of the WriteError,
// so that they appear together in the log output.
func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("err", e.Err),
		slog.String("code", e.Code),
		slog.String("message", e.Message),
		slog.String("kind", fmt.Sprintf("%T", e)))
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	} else if errors.Is(err, context.Canceled) {
		return ECANCELED
	}
	return EINTERNAL
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	var e *Error
	if errors.As(err, &e) {
		return e.Message
	}

	return "Internal error."
}

// Errorf is a helper function to return an WriteError with a given code and formatted message.
func Errorf(code, msgWithFormat string, args ...any) *Error {
	if len(args) > 0 {
		msgWithFormat = fmt.Sprintf(msgWithFormat, args...)
	}
	return &Error{
		Code:    code,
		Message: msgWithFormat,
	}
}

// ErrorStatusCode returns the associated HTTP status code for an arh.WriteError code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated arh.code for a HTTP status code.
func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return EINTERNAL
}
