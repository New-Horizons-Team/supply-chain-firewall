package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type LogLevel string

const (
	LevelDebug   LogLevel = "DEBUG"
	LevelInfo    LogLevel = "INFO"
	LevelWarning LogLevel = "WARNING"
	LevelError   LogLevel = "ERROR"
)

const DefaultLevel = LevelWarning

// GetLevels retrieves the available levels, properly converted to strings.
func GetLevels() []string {
	return []string{
		string(LevelDebug),
		string(LevelInfo),
		string(LevelWarning),
		string(LevelError),
	}
}

var currentLevel = DefaultLevel

// GetCurrentLevel retrieves the current log level.
func GetCurrentLevel() LogLevel {
	return currentLevel
}

// Configure configures the global logger.
func Configure(level LogLevel) {
	logrus.SetFormatter(&formatter{})

	var logrusLevel logrus.Level

	if level == "" {
		level = DefaultLevel
	}

	switch level {
	case LevelDebug:
		logrusLevel = logrus.DebugLevel
	case LevelInfo:
		logrusLevel = logrus.InfoLevel
	case LevelWarning:
		logrusLevel = logrus.WarnLevel
	case LevelError:
		logrusLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logrusLevel)
	currentLevel = level
}

// Default formatter.
type formatter struct{}

// Format implements logrus.Formatter.
func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var level LogLevel

	switch entry.Level {
	case logrus.DebugLevel:
		level = LevelDebug
	case logrus.InfoLevel:
		level = LevelInfo
	case logrus.WarnLevel:
		level = LevelWarning
	case logrus.ErrorLevel:
		level = LevelError
	}

	msg := fmt.Sprintf("[SCFW] %s: %s\n", level, entry.Message)

	return []byte(msg), nil
}

type DefaultLogger struct{}

// Debug implements pkg/logger.Logger.
func (DefaultLogger) Debug(ctx context.Context, format string, args ...any) {
	logrus.Debugf(format, args...)
}

// Info implements pkg/logger.Logger.
func (DefaultLogger) Info(ctx context.Context, format string, args ...any) {
	logrus.Infof(format, args...)
}

// Warn implements pkg/logger.Logger.
func (DefaultLogger) Warn(ctx context.Context, format string, args ...any) {
	logrus.Warnf(format, args...)
}

// Error implements pkg/logger.Logger.
func (DefaultLogger) Error(ctx context.Context, format string, args ...any) {
	logrus.Errorf(format, args...)
}

type mutedLogger struct{}

// Debug implements pkg/logger.Logger.
func (mutedLogger) Debug(ctx context.Context, format string, args ...any) {
}

// Info implements pkg/logger.Logger.
func (mutedLogger) Info(ctx context.Context, format string, args ...any) {
}

// Warn implements pkg/logger.Logger.
func (mutedLogger) Warn(ctx context.Context, format string, args ...any) {
}

// Error implements pkg/logger.Logger.
func (mutedLogger) Error(ctx context.Context, format string, args ...any) {
}
