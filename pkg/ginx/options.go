package ginx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"research-apm/pkg/ginx/internal/auth"
	"research-apm/pkg/ginx/internal/logger"
	"research-apm/pkg/ginx/internal/traceid"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
)

// WithVerifyHMAC adds a middleware that verifies HMAC signatures
// on incoming requests using the provided secret key.
func WithVerifyHMAC(secret string) EngineOption {
	return func(e *gin.Engine) {
		e.Use(auth.VerifyHMAC(secret))
	}
}

// LogConfig holds metadata information to be included in log entries.
type LogConfig struct {
	AppName      string
	AppSite      string
	AppEnv       string
	AppVersion   string
	AppDBVersion string
}

// WithLogFile adds a middleware that logs request/response information
// and writes logs into the given file in JSON format.
//
// Parameters:
//   - ctx: context for managing goroutine lifecycle (used for graceful shutdown)
//   - file: destination writer (e.g., os.File)
//   - config: application metadata included in each log entry
func WithLogFile(
	ctx context.Context,
	file io.Writer,
	config LogConfig,
) EngineOption {
	encoder := json.NewEncoder(file)
	logCh := make(chan logger.Logging, 100)

	var wg sync.WaitGroup
	wg.Add(1)

	// Background goroutine that encodes logs to the file
	go func() {
		defer wg.Done()
		for {
			select {
			case log, ok := <-logCh:
				if !ok {
					return // channel closed
				}
				encoder.Encode(log)
			case <-ctx.Done():
				return // shutdown signal
			}
		}
	}()

	// Shutdown hook: wait for all logs to be flushed before exit
	go func() {
		<-ctx.Done()
		close(logCh)
		wg.Wait()
	}()

	return func(e *gin.Engine) {
		e.Use(logger.NewLogger(
			config.AppName,
			config.AppSite,
			config.AppEnv,
			config.AppVersion,
			config.AppDBVersion,
			logCh,
		))
	}
}

// WithLogPushHttp adds a middleware that logs request/response information
// and pushes logs to a remote HTTP endpoint in JSON format.
//
// Parameters:
//   - ctx: context for managing goroutine lifecycle (used for graceful shutdown)
//   - url: target endpoint where logs will be pushed
//   - header: custom headers to attach (e.g., Authorization tokens)
//   - config: application metadata included in each log entry
func WithLogPushHttp(
	ctx context.Context,
	url string,
	header map[string]string,
	config LogConfig,
) EngineOption {
	logCh := make(chan logger.Logging, 100)

	var wg sync.WaitGroup
	wg.Add(1)

	// Background goroutine that pushes logs to the HTTP endpoint
	go func() {
		defer wg.Done()
		client := &http.Client{Timeout: 5 * time.Second}

		for {
			select {
			case log, ok := <-logCh:
				if !ok {
					return // channel closed
				}
				data, err := json.Marshal(log)
				if err != nil {
					fmt.Println("[ERROR] log push http: failed to marshal:", err)
					continue
				}

				req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
				if err != nil {
					fmt.Println("[ERROR] log push http: failed to create request:", err)
					continue
				}

				req.Header.Set("Content-Type", "application/json")
				for k, v := range header {
					req.Header.Set(k, v)
				}

				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("[ERROR] log push http: failed to send:", err)
					continue
				}
				resp.Body.Close()

			case <-ctx.Done():
				return // shutdown signal
			}
		}
	}()

	// Shutdown hook: wait for all logs to be flushed before exit
	go func() {
		<-ctx.Done()
		close(logCh)
		wg.Wait()
	}()

	return func(e *gin.Engine) {
		e.Use(logger.NewLogger(
			config.AppName,
			config.AppSite,
			config.AppEnv,
			config.AppVersion,
			config.AppDBVersion,
			logCh,
		))
	}
}

// WithTraceID adds a middleware that ensures every request has a trace ID (X-Trace-ID).
// If the request does not include a trace ID, a new one is generated.
// The trace ID is then attached to both request and response headers,
// making it useful for distributed tracing and debugging.
func WithTraceID() EngineOption {
	return func(e *gin.Engine) {
		e.Use(traceid.TraceID())
	}
}

// WithElasticAPM adds Elastic APM middleware to the Gin engine.
// It automatically instruments incoming HTTP requests for performance
// monitoring and error tracking. Requires proper Elastic APM configuration
// via environment variables (e.g., ELASTIC_APM_SERVER_URL, ELASTIC_APM_SERVICE_NAME).
func WithElasticAPM() EngineOption {
	return func(e *gin.Engine) {
		e.Use(apmgin.Middleware(e))
	}
}

// WithCors adds CORS (Cross-Origin Resource Sharing) middleware to the Gin engine.
// This enables the server to handle cross-origin requests based on the provided config.
//
// Example usage:
//
//	ginx.WithCors(cors.Config{
//	    AllowOrigins: []string{"https://example.com"},
//	    AllowMethods: []string{"GET", "POST"},
//	})
func WithCors(config cors.Config) EngineOption {
	return func(e *gin.Engine) {
		e.Use(cors.New(config))
	}
}
