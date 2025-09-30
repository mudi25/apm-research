package mongox

import (
	"context"
	"errors"
	appErr "research-apm/pkg/errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// isRetryable determines whether a MongoDB error is transient and safe to retry.
func isRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Network or timeout errors are considered retryable.
	// These errors usually originate from the transport layer (TCP/socket)
	// and do not contain MongoDB-specific error codes.
	if mongo.IsNetworkError(err) || mongo.IsTimeout(err) {
		return true
	}

	// List of MongoDB server error codes considered retryable.
	// Source: https://www.mongodb.com/docs/manual/core/retryable-writes/#retryable-errors
	var retryableCodes = map[int]bool{
		6:     true, // HostUnreachable
		7:     true, // HostNotFound
		89:    true, // NetworkTimeout
		91:    true, // ShutdownInProgress
		189:   true, // PrimarySteppedDown
		262:   true, // ExceededTimeLimit
		9001:  true, // SocketException
		10107: true, // NotMaster (deprecated, still relevant)
		11600: true, // InterruptedAtShutdown
		11602: true, // InterruptedDueToReplStateChange
		13435: true, // NotPrimaryNoSecondaryOk
		13436: true, // NotPrimaryOrSecondary
	}

	// Check CommandError (typically from read/write commands like find, aggregate, etc.)
	var cmdErr mongo.CommandError
	if errors.As(err, &cmdErr) {
		if retryableCodes[int(cmdErr.Code)] {
			return true
		}
	}

	// Check WriteException (from insert, update, delete operations)
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		// Retry if any individual write error has a retryable code
		for _, we := range writeErr.WriteErrors {
			if retryableCodes[we.Code] {
				return true
			}
		}
		// Retry if the write concern error is retryable
		if wcErr := writeErr.WriteConcernError; wcErr != nil {
			if retryableCodes[wcErr.Code] {
				return true
			}
		}
	}

	return false
}

func NewError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if isRetryable(err) {
		err = appErr.NewRetryable(err)
	}
	return err
}
