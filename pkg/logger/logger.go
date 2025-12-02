package logger

import "context"

// Logger provides a unified interface for logging messages for the user.
// Any object that implements Logger is responsible for filtering
// messages in an undesired log level.
//
// IMPORTANT: Any object that implements Logger SHALL be thread safe!
type Logger interface {
	// Debug logs a debug message.
	Debug(ctx context.Context, format string, args ...any)

	// Info logs an info message.
	Info(ctx context.Context, format string, args ...any)

	// Warn logs a warning message.
	Warn(ctx context.Context, format string, args ...any)

	// Error logs an error message.
	Error(ctx context.Context, format string, args ...any)
}
