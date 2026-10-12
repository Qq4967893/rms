package bootstrap

import (
	"backend-rms/core/config"
	"backend-rms/core/logger"

	"github.com/getsentry/sentry-go"
)

func InitSentry() {
	cfg := config.Conf.Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Dsn,
		EnableTracing:    cfg.EnableTracing,
		TracesSampleRate: cfg.TracesSampleRate,
	}); err != nil {
		logger.Errorf("Sentry initialization failed: %v\n", err)
	}
}
