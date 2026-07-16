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
//
// The path is resolved through symlinks so that it compares equal to os.Getwd()
// after Chdir: on macOS t.TempDir() returns /var/..., but /var is a symlink to
// /private/var and Getwd reports the resolved form.
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("failed to resolve temp directory: %v", err)
	}
	return dir
}

// TempFile creates a temporary file with the given content.
// Returns the file path. The file is automatically cleaned up.
func TempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := TempDir(t)
	path := filepath.Join(dir, name)
	// #nosec G306 -- test helper writes non-secret fixtures; 0644 is intentional
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
	// #nosec G301 -- test helper creates temp dirs; 0755 is intentional
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	// #nosec G306 -- test helper writes non-secret fixtures; 0644 is intentional
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

	rOut, wOut, err := os.Pipe()
	if err != nil {
		panic("testutil.Capture: failed to create stdout pipe: " + err.Error())
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		panic("testutil.Capture: failed to create stderr pipe: " + err.Error())
	}

	os.Stdout = wOut
	os.Stderr = wErr

	fn()

	// #nosec G104 -- close errors on pipe writers are non-actionable after capture
	_ = wOut.Close()
	// #nosec G104 -- close errors on pipe writers are non-actionable after capture
	_ = wErr.Close()

	var bufOut, bufErr bytes.Buffer
	if _, err := io.Copy(&bufOut, rOut); err != nil {
		panic("testutil.Capture: failed to read stdout: " + err.Error())
	}
	if _, err := io.Copy(&bufErr, rErr); err != nil {
		panic("testutil.Capture: failed to read stderr: " + err.Error())
	}

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
	r, w, err := os.Pipe()
	if err != nil {
		panic("testutil.CaptureStdout: failed to create pipe: " + err.Error())
	}
	os.Stdout = w

	fn()

	// #nosec G104 -- close errors on pipe writers are non-actionable after capture
	_ = w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic("testutil.CaptureStdout: failed to read stdout: " + err.Error())
	}
	os.Stdout = oldStdout

	return buf.String()
}

// CaptureStderr captures only stderr during the execution of fn.
func CaptureStderr(fn func()) string {
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic("testutil.CaptureStderr: failed to create pipe: " + err.Error())
	}
	os.Stderr = w

	fn()

	// #nosec G104 -- close errors on pipe writers are non-actionable after capture
	_ = w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic("testutil.CaptureStderr: failed to read stderr: " + err.Error())
	}
	os.Stderr = oldStderr

	return buf.String()
}
