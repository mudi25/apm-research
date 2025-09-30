package response

import (
	"fmt"
	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/tracer"

	"github.com/gin-gonic/gin"
)

// Response defines the standard API response format.
type Response struct {
	Code        string  `json:"code"`        // Application-level status code
	Message     string  `json:"message"`     // Human-readable status message
	Data        any     `json:"data"`        // Payload data (if any)
	Errors      *string `json:"errors"`      // Optional detailed error message
	IsRetryable bool    `json:"isRetryable"` // Indicates whether the error can be retried
}

// New sends a standard JSON response based on the result or the provided error.
// If an error is present, it maps it to an error response using the custom error package.
func New(ctx *gin.Context, result any, err error) {
	if err := errors.FromError(err); err != nil {
		merr := err.Error()
		tracer.CaptureError(ctx.Request.Context(), err)
		ctx.JSON(err.Code.HttpStatus(), &Response{
			Code:        string(err.Code),
			Message:     err.Message,
			Data:        nil,
			Errors:      &merr,
			IsRetryable: err.IsRetryable,
		})
		return
	}

	// Success response
	ctx.JSON(codes.Success.HttpStatus(), &Response{
		Code:        string(codes.Success),
		Message:     "success",
		Data:        result,
		Errors:      nil,
		IsRetryable: false,
	})
}

// Abort is the same as New, but also stops the middleware chain by calling AbortWithStatusJSON.
// Use this when you want to return early and prevent further processing.
func Abort(ctx *gin.Context, err error) {
	merr := errors.New(codes.UnknownError, "request abort with unknown error", fmt.Errorf("unknown error"))
	if err := errors.FromError(err); err != nil {
		merr = err
	}
	tracer.CaptureError(ctx.Request.Context(), merr)
	ctx.AbortWithStatusJSON(merr.Code.HttpStatus(), &Response{
		Code:        string(merr.Code),
		Message:     merr.Error(),
		Data:        nil,
		Errors:      &merr.Message,
		IsRetryable: false,
	})
}
