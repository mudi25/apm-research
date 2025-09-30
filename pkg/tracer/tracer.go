package tracer

import (
	"context"
	"os"

	"go.elastic.co/apm/v2"
)

type Config struct {
	Env            string
	ServiceName    string
	Version        string
	ServerUrl      string
	SecretToken    string
	IsUsingLogging bool
}

func InitTracer(cfg Config) error {
	os.Setenv("ELASTIC_APM_ENVIRONMENT", cfg.Env)
	os.Setenv("ELASTIC_APM_SERVICE_VERSION", cfg.Version)
	os.Setenv("ELASTIC_APM_SERVICE_NAME", cfg.ServiceName)
	os.Setenv("ELASTIC_APM_SERVER_URL", cfg.ServerUrl)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", cfg.SecretToken)
	if cfg.IsUsingLogging {
		os.Setenv("ELASTIC_APM_LOG_LEVEL", "debug")
		os.Setenv("ELASTIC_APM_LOG_FILE", "stdout")
	}
	return nil
}

func StartSpan(ctx context.Context, name string) (context.Context, *apm.Span) {
	tx := apm.TransactionFromContext(ctx)
	if tx == nil {
		tx = apm.DefaultTracer().StartTransaction(name, "custom")
	}
	span, ctx := apm.StartSpan(ctx, name, tx.Type)
	return ctx, span
}

func CaptureError(ctx context.Context, err error) {
	apm.CaptureError(ctx, err).Send()
}
