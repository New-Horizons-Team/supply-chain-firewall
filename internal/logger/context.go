package logger

import (
	"context"

	"github.com/New-Horizons-Team/supply-chain-firewall/pkg/logger"
)

type loggerKey struct{}

// WithLogger adds the logger to the context.
func WithLogger(ctx context.Context, logger logger.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLogger retrieves a logger from the context,
// returning the muted logger if none is found.
func GetLogger(ctx context.Context) logger.Logger {
	var defaultLogger mutedLogger

	logger, _ := ctx.Value(loggerKey{}).(logger.Logger)
	if logger != nil {
		return logger
	}

	return defaultLogger
}
