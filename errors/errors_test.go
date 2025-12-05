package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	inner := New("inner error")
	outer := New("outer error")

	wrapped := Wrap(inner, outer)
	if wrapped == nil {
		t.Fatal("expected non-nil error")
	}

	// Should match both errors
	if !Is(wrapped, inner) {
		t.Error("wrapped error should match inner")
	}
	if !Is(wrapped, outer) {
		t.Error("wrapped error should match outer")
	}
}

func TestWrap_NilHandling(t *testing.T) {
	err := New("test error")

	if got := Wrap(nil, err); got != err {
		t.Errorf("Wrap(nil, err) = %v, want %v", got, err)
	}
	if got := Wrap(err, nil); got != err {
		t.Errorf("Wrap(err, nil) = %v, want %v", got, err)
	}
	if got := Wrap(nil, nil); got != nil {
		t.Errorf("Wrap(nil, nil) = %v, want nil", got)
	}
}

func TestWrapWithMessage(t *testing.T) {
	err := New("original error")
	wrapped := WrapWithMessage(err, "context message")

	if !strings.Contains(wrapped.Error(), "context message") {
		t.Error("expected context message in error")
	}
	if !strings.Contains(wrapped.Error(), "original error") {
		t.Error("expected original error in wrapped error")
	}
}

func TestWrapWithMessage_Nil(t *testing.T) {
	if got := WrapWithMessage(nil, "message"); got != nil {
		t.Errorf("WrapWithMessage(nil, msg) = %v, want nil", got)
	}
}

func TestWrapOp(t *testing.T) {
	err := New("file not accessible")
	wrapped := WrapOp("open file", err)

	if !strings.Contains(wrapped.Error(), "open file failed") {
		t.Errorf("expected 'open file failed' in error, got: %s", wrapped.Error())
	}
}

func TestWrapOp_Nil(t *testing.T) {
	if got := WrapOp("operation", nil); got != nil {
		t.Errorf("WrapOp(op, nil) = %v, want nil", got)
	}
}

func TestNew(t *testing.T) {
	err := New("test error")
	if err.Error() != "test error" {
		t.Errorf("got %q, want %q", err.Error(), "test error")
	}
}

func TestNewf(t *testing.T) {
	err := Newf("error %d: %s", 42, "test")
	if err.Error() != "error 42: test" {
		t.Errorf("got %q, want %q", err.Error(), "error 42: test")
	}
}

func TestIs(t *testing.T) {
	if !Is(ErrNotFound, ErrNotFound) {
		t.Error("ErrNotFound should match itself")
	}

	wrapped := WrapWithMessage(ErrNotFound, "context")
	if !Is(wrapped, ErrNotFound) {
		t.Error("wrapped error should match sentinel")
	}
}

type customError struct {
	Code int
}

func (e *customError) Error() string {
	return fmt.Sprintf("error code: %d", e.Code)
}

func TestAs(t *testing.T) {
	err := &customError{Code: 42}
	wrapped := WrapWithMessage(err, "context")

	var target *customError
	if !As(wrapped, &target) {
		t.Error("As should find customError")
	}
	if target.Code != 42 {
		t.Errorf("got code %d, want 42", target.Code)
	}
}

func TestJoin(t *testing.T) {
	err1 := New("error 1")
	err2 := New("error 2")
	joined := Join(err1, err2)

	if !Is(joined, err1) {
		t.Error("joined error should match err1")
	}
	if !Is(joined, err2) {
		t.Error("joined error should match err2")
	}
}

func TestUnwrap(t *testing.T) {
	inner := New("inner")
	wrapped := WrapWithMessage(inner, "outer")

	unwrapped := Unwrap(wrapped)
	if unwrapped == nil {
		t.Error("expected unwrapped error")
	}
}

func TestSentinelErrors(t *testing.T) {
	sentinels := []error{
		ErrNotFound,
		ErrInvalidInput,
		ErrConfigNotFound,
		ErrInvalidConfig,
		ErrUnauthorized,
		ErrTimeout,
		ErrPermission,
		ErrAlreadyExists,
		ErrNotSupported,
	}

	for _, err := range sentinels {
		if err == nil {
			t.Error("sentinel error should not be nil")
		}
		if err.Error() == "" {
			t.Error("sentinel error should have message")
		}
	}
}

// Validation error tests

func TestInvalidPath(t *testing.T) {
	err := InvalidPath("config", errors.New("not readable"))
	if !strings.Contains(err.Error(), "invalid config path") {
		t.Errorf("unexpected error: %s", err.Error())
	}

	err2 := InvalidPath("source", nil)
	if !strings.Contains(err2.Error(), "invalid source path") {
		t.Errorf("unexpected error: %s", err2.Error())
	}
}

func TestFileNotFound(t *testing.T) {
	err := FileNotFound("/path/to/file")
	if !strings.Contains(err.Error(), "file not found") {
		t.Error("expected 'file not found' in error")
	}
	if !strings.Contains(err.Error(), "/path/to/file") {
		t.Error("expected path in error")
	}
}

func TestDirNotFound(t *testing.T) {
	err := DirNotFound("/path/to/dir")
	if !strings.Contains(err.Error(), "directory does not exist") {
		t.Error("expected 'directory does not exist' in error")
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError("invalid format")
	if !strings.Contains(err.Error(), "validation error") {
		t.Error("expected 'validation error' prefix")
	}
}

func TestValidationErrorf(t *testing.T) {
	err := ValidationErrorf("field %s is invalid", "name")
	if !strings.Contains(err.Error(), "field name is invalid") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRequiredFlag(t *testing.T) {
	err := RequiredFlag("output")
	if !strings.Contains(err.Error(), "--output flag is required") {
		t.Error("expected flag name in error")
	}

	errWithExamples := RequiredFlag("format", "gz-tool --format=json", "gz-tool --format=yaml")
	if !strings.Contains(errWithExamples.Error(), "Examples:") {
		t.Error("expected examples in error")
	}
}

func TestMutuallyExclusive(t *testing.T) {
	err := MutuallyExclusive("verbose", "quiet")
	if !strings.Contains(err.Error(), "--verbose") {
		t.Error("expected first flag in error")
	}
	if !strings.Contains(err.Error(), "--quiet") {
		t.Error("expected second flag in error")
	}
	if !strings.Contains(err.Error(), "cannot be used together") {
		t.Error("expected 'cannot be used together' in error")
	}
}

func TestMinValue(t *testing.T) {
	err := MinValue("count", 1)
	if !strings.Contains(err.Error(), "must be at least 1") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestMaxValue(t *testing.T) {
	err := MaxValue("limit", 100)
	if !strings.Contains(err.Error(), "must be at most 100") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestRange(t *testing.T) {
	err := Range("port", 1, 65535)
	if !strings.Contains(err.Error(), "must be between 1 and 65535") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestEmptyValue(t *testing.T) {
	err := EmptyValue("name")
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestInvalidValue(t *testing.T) {
	err := InvalidValue("format", "xyz", "must be json or yaml")
	if !strings.Contains(err.Error(), "invalid format") {
		t.Error("expected field name in error")
	}
	if !strings.Contains(err.Error(), "xyz") {
		t.Error("expected value in error")
	}

	err2 := InvalidValue("format", "xyz", "")
	if !strings.Contains(err2.Error(), "invalid format") {
		t.Error("expected field name in error without reason")
	}
}
