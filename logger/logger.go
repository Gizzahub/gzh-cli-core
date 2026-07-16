// Package logger provides structured logging utilities for gzh-cli tools.
package logger

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
	"sync"
	"time"
)

// Level represents logging level.
type Level int

// Log levels from most to least verbose.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of the level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a string into a Level.
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// Logger defines the logging interface.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	WithContext(key string, value any) Logger
	SetLevel(level Level)
	SetOutput(w io.Writer)
}

// SimpleLogger is a basic structured logger implementation.
type SimpleLogger struct {
	mu      sync.Mutex
	name    string
	level   Level
	out     io.Writer
	context map[string]any
}

// New creates a new SimpleLogger with the given name.
func New(name string) *SimpleLogger {
	return &SimpleLogger{
		name:    name,
		level:   LevelInfo,
		out:     os.Stdout,
		context: make(map[string]any),
	}
}

// SetOutput sets the output destination.
func (l *SimpleLogger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetLevel sets the minimum logging level.
func (l *SimpleLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current logging level.
func (l *SimpleLogger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// WithContext returns a new logger with additional context.
func (l *SimpleLogger) WithContext(key string, value any) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newContext := make(map[string]any)
	maps.Copy(newContext, l.context)
	newContext[key] = value

	return &SimpleLogger{
		name:    l.name,
		level:   l.level,
		out:     l.out,
		context: newContext,
	}
}

func (l *SimpleLogger) log(level Level, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Build context string.
	var contextStr strings.Builder
	for k, v := range l.context {
		fmt.Fprintf(&contextStr, " %s=%v", k, v)
	}

	// Build args string (key=value pairs).
	var argsStr strings.Builder
	for i := 0; i < len(args)-1; i += 2 {
		if i+1 < len(args) {
			fmt.Fprintf(&argsStr, " %v=%v", args[i], args[i+1])
		}
	}

	fmt.Fprintf(l.out, "[%s] %s [%s]%s %s%s\n",
		timestamp, level.String(), l.name, contextStr.String(), msg, argsStr.String())
}

// Debug logs a debug message.
func (l *SimpleLogger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args...)
}

// Info logs an info message.
func (l *SimpleLogger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args...)
}

// Warn logs a warning message.
func (l *SimpleLogger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args...)
}

// Error logs an error message.
func (l *SimpleLogger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args...)
}

// NopLogger is a no-operation logger that discards all output.
type NopLogger struct{}

// NewNop creates a new no-operation logger.
func NewNop() *NopLogger {
	return &NopLogger{}
}

// Debug discards the debug message.
func (l *NopLogger) Debug(msg string, args ...any) {}

// Info discards the info message.
func (l *NopLogger) Info(msg string, args ...any) {}

// Warn discards the warning message.
func (l *NopLogger) Warn(msg string, args ...any) {}

// Error discards the error message.
func (l *NopLogger) Error(msg string, args ...any) {}

// WithContext returns the same no-op logger.
func (l *NopLogger) WithContext(key string, value any) Logger {
	return l
}

// SetLevel is a no-op.
func (l *NopLogger) SetLevel(level Level) {}

// SetOutput is a no-op.
func (l *NopLogger) SetOutput(w io.Writer) {}

// Ensure interfaces are implemented.
var (
	_ Logger = (*SimpleLogger)(nil)
	_ Logger = (*NopLogger)(nil)
)
