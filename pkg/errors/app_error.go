package errors

import (
	"fmt"
	"research-apm/pkg/errors/codes"
)

// AppError is a structured application-level error
// that includes a code, message, the original error, and retryability flag.
type AppError struct {
	Code        codes.Code // custom application error code
	Message     string     // user-friendly message
	Errors      error      // underlying error
	IsRetryable bool       // indicates if the error is safe to retry
}

// Error implements the error interface for AppError.
// If the underlying error is nil, it returns a composed fallback message.
func (e *AppError) Error() string {
	if e.Errors != nil {
		return e.Errors.Error()
	}
	return fmt.Sprintf("apperror is nil with code %s and message %s | retryable: %t", e.Code, e.Message, e.IsRetryable)
}

// Wrap wraps a standard error into an AppError with optional retry logic.
// - If the error is nil, it returns an UnknownError.
// - If the error is a Retryable, it converts it to an AppError with retryable = true.
// - If the error is already an AppError, it returns it as-is.
func Wrap(code codes.Code, message string, err error) *AppError {
	if err == nil {
		return &AppError{Code: codes.UnknownError, Message: message, Errors: fmt.Errorf("error is nil"), IsRetryable: false}
	}
	retry := false
	if e, ok := err.(*Retryable); ok && e.IsRetryable {
		message = "The system is busy, please try again later."
		code = codes.Unavailable
		retry = true
	} else if e, ok := err.(*AppError); ok {
		return e
	}
	return &AppError{Code: code, Message: message, Errors: err, IsRetryable: retry}
}

// FromError converts any error into an AppError.
// If it's already an AppError, it returns it unchanged.
// Otherwise, it wraps it with an UnknownError code.
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return &AppError{
		Code:        codes.UnknownError,
		Message:     "unknown error",
		Errors:      err,
		IsRetryable: false,
	}
}

// New creates a new AppError with the given code, message, and underlying error.
func New(code codes.Code, message string, err error) *AppError {
	return &AppError{
		Code:        code,
		Message:     message,
		Errors:      err,
		IsRetryable: false,
	}
}

// NewBadRequest creates an AppError with a BadRequest code.
func NewBadRequest(message string, err error) *AppError {
	return &AppError{
		Code:        codes.BadRequest,
		Message:     message,
		Errors:      err,
		IsRetryable: false,
	}
}
