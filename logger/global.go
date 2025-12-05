// Package logger provides structured logging utilities for gzh-cli tools.
package logger

import "io"

// defaultLogger is the global default logger instance.
var defaultLogger Logger = New("app")

// SetDefault sets the default global logger.
func SetDefault(l Logger) {
	defaultLogger = l
}

// Default returns the default global logger.
func Default() Logger {
	return defaultLogger
}

// SetDefaultLevel sets the level of the default logger.
func SetDefaultLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetDefaultOutput sets the output of the default logger.
func SetDefaultOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// Debug logs a debug message using the default logger.
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message using the default logger.
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger.
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}
