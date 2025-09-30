package response_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/ginx/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// init runs once before all tests to set Gin to test mode
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

// TestNewSuccess verifies that response.New returns a successful JSON response
// with status 200 and the correct payload when no error is passed.
func TestNewSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	data := map[string]string{"hello": "world"}

	response.New(ctx, data, nil)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"code": "SUCCESS",
		"message": "success",
		"data": {"hello": "world"},
		"errors": null,
		"isRetryable": false
	}`, w.Body.String())
}

// TestNewWithError verifies that response.New returns an error response
// when a custom application error is provided.
func TestNewWithError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := errors.New(codes.BadRequest, "Invalid request", fmt.Errorf("invalid request"))

	response.New(ctx, nil, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"BAD_REQUEST"`)
	assert.Contains(t, w.Body.String(), `"message":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"isRetryable":false`)
}

// TestAbortSuccess checks that response.Abort returns a success response
// and stops further middleware execution.
func TestAbortUnknownError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	response.Abort(ctx, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{
		"code": "UNKNOWN_ERROR",
		"message": "request abort with unknown error",
		"data": null,
		"errors": "request abort with unknown error",
		"isRetryable": false
	}`, w.Body.String())
}

// TestAbortWithError verifies that response.Abort returns the correct
// error response and stops further processing when an error is passed.
func TestAbortWithError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := errors.New(codes.PermissionDenied, "permission denied", fmt.Errorf("permission denied"))

	response.Abort(ctx, err)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{
		"code": "PERMISSION_DENIED",
		"message": "permission denied",
		"data": null,
		"errors": "permission denied",
		"isRetryable": false
	}`, w.Body.String())
}

// TestNewWithStandardError ensures response.New can handle a standard Go error
// (not wrapped in the custom error package) and returns a generic internal error.
func TestNewWithStandardError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	err := fmt.Errorf("standard error")

	response.New(ctx, nil, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"code":"UNKNOWN_ERROR"`)
	assert.Contains(t, w.Body.String(), `"errors":"standard error"`)
}
