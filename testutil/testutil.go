// Package testutil provides testing utilities and helpers for gzh-cli tools.
package testutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory and returns its path.
// The directory is automatically cleaned up when the test finishes.
func TempDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// TempFile creates a temporary file with the given content.
// Returns the file path. The file is automatically cleaned up.
func TempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return path
}

// TempFileInDir creates a temporary file in the specified directory.
// Returns the file path.
func TempFileInDir(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return path
}

// CaptureOutput holds captured stdout and stderr output.
type CaptureOutput struct {
	Stdout string
	Stderr string
}

// Capture captures stdout and stderr during the execution of fn.
func Capture(fn func()) CaptureOutput {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	fn()

	wOut.Close()
	wErr.Close()

	var bufOut, bufErr bytes.Buffer
	io.Copy(&bufOut, rOut)
	io.Copy(&bufErr, rErr)

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return CaptureOutput{
		Stdout: bufOut.String(),
		Stderr: bufErr.String(),
	}
}

// CaptureStdout captures only stdout during the execution of fn.
func CaptureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldStdout

	return buf.String()
}

// CaptureStderr captures only stderr during the execution of fn.
func CaptureStderr(fn func()) string {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stderr = oldStderr

	return buf.String()
}
