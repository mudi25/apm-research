package errors

// Retryable wraps an error and marks it as retryable.
// Useful for categorizing transient errors (e.g. timeouts, temporary service issues).
type Retryable struct {
	Errors      error
	IsRetryable bool
}

// Error implements the error interface for Retryable.
// If the wrapped error is nil, it returns a default message.
func (e *Retryable) Error() string {
	if e.Errors != nil {
		return e.Errors.Error()
	}
	return "retryable error is nil"
}

// NewRetryable creates a new retryable error from a given error.
func NewRetryable(err error) *Retryable {
	return &Retryable{Errors: err, IsRetryable: true}
}

// check is error is retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.IsRetryable
	}
	if appErr, ok := err.(*Retryable); ok {
		return appErr.IsRetryable
	}
	return false
}
