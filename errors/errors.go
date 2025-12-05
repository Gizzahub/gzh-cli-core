// Package errors provides common error types and utilities for gzh-cli tools.
package errors

import (
	"errors"
	"fmt"
)

// Standard sentinel errors - common across all gzh-cli projects.
var (
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrConfigNotFound = errors.New("config not found")
	ErrInvalidConfig  = errors.New("invalid config")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrTimeout        = errors.New("operation timed out")
	ErrPermission     = errors.New("permission denied")
	ErrAlreadyExists  = errors.New("already exists")
	ErrNotSupported   = errors.New("not supported")
)

// Wrap combines two errors, preserving the chain for errors.Is/As.
// If err is nil, returns target. If target is nil, returns err.
func Wrap(err, target error) error {
	if err == nil {
		return target
	}
	if target == nil {
		return err
	}
	return fmt.Errorf("%w: %w", target, err)
}

// WrapWithMessage wraps an error with additional context message.
func WrapWithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapOp wraps an error with operation context.
func WrapOp(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s failed: %w", operation, err)
}

// New creates a new error with the given message.
func New(message string) error {
	return errors.New(message)
}

// Newf creates a new error with the formatted message.
func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// Is reports whether any error in err's tree matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's tree that matches target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Join returns an error that wraps the given errors.
func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Unwrap returns the result of calling the Unwrap method on err.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
