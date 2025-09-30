package ginx

import (
	"fmt"
	"os"
	"strings"

	"research-apm/pkg/errors"
	"research-apm/pkg/errors/codes"
	"research-apm/pkg/ginx/response"

	"github.com/gin-gonic/gin"
)

type EngineOption func(*gin.Engine)

// NewEngine creates and configures a new Gin engine instance.
// It performs the following:
//
// 1. Checks the ENV environment variable:
//   - If ENV starts with "prod" (case-insensitive), Gin runs in ReleaseMode
//     which disables debug logs and reduces console output.
//   - Otherwise, Gin runs in its default development mode.
//
// 2. Sets a fallback handler for unmatched routes returning a 404 response.
func NewEngine(options ...EngineOption) *gin.Engine {
	if strings.HasPrefix(strings.ToLower(os.Getenv("ENV")), "prod") {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	for _, opt := range options {
		opt(engine)
	}

	engine.NoRoute(func(ctx *gin.Context) {
		response.New(ctx, nil, errors.New(
			codes.PathNotFound,
			"request path not found",
			fmt.Errorf("request path not found"),
		))
	})

	return engine
}
