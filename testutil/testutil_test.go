package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestTempDir(t *testing.T) {
	dir := TempDir(t)
	if dir == "" {
		t.Fatal("expected non-empty directory path")
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestTempFile(t *testing.T) {
	content := "test content"
	path := TempFile(t, "test.txt", content)

	if path == "" {
		t.Fatal("expected non-empty path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != content {
		t.Errorf("got %q, want %q", string(data), content)
	}
}

func TestTempFileInDir(t *testing.T) {
	dir := TempDir(t)
	content := "nested content"
	path := TempFileInDir(t, dir, "subdir/test.txt", content)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != content {
		t.Errorf("got %q, want %q", string(data), content)
	}
}

func TestCapture(t *testing.T) {
	output := Capture(func() {
		fmt.Fprintln(os.Stdout, "stdout message")
		fmt.Fprintln(os.Stderr, "stderr message")
	})

	AssertContains(t, output.Stdout, "stdout message")
	AssertContains(t, output.Stderr, "stderr message")
}

func TestCaptureStdout(t *testing.T) {
	output := CaptureStdout(func() {
		fmt.Println("stdout only")
	})
	AssertContains(t, output, "stdout only")
}

func TestCaptureStderr(t *testing.T) {
	output := CaptureStderr(func() {
		fmt.Fprintln(os.Stderr, "stderr only")
	})
	AssertContains(t, output, "stderr only")
}

func TestAssertNoError(t *testing.T) {
	// Should not fail
	AssertNoError(t, nil)
}

func TestAssertError(t *testing.T) {
	// Should not fail
	AssertError(t, errors.New("some error"))
}

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, 1, 1)
	AssertEqual(t, "hello", "hello")
	AssertEqual(t, []int{1, 2, 3}, []int{1, 2, 3})
}

func TestAssertContains(t *testing.T) {
	AssertContains(t, "hello world", "world")
}

func TestAssertNotContains(t *testing.T) {
	AssertNotContains(t, "hello world", "foo")
}

func TestSetEnv(t *testing.T) {
	key := "GZH_TEST_ENV_VAR"
	original := os.Getenv(key)
	defer os.Setenv(key, original)

	SetEnv(t, key, "test_value")
	if got := os.Getenv(key); got != "test_value" {
		t.Errorf("got %q, want %q", got, "test_value")
	}
}

func TestChdir(t *testing.T) {
	original, _ := os.Getwd()
	defer os.Chdir(original) // ensure we restore even if test fails

	dir := TempDir(t)
	Chdir(t, dir)

	current, _ := os.Getwd()
	if current != dir {
		t.Errorf("got %q, want %q", current, dir)
	}
}

func TestChdirTemp(t *testing.T) {
	original, _ := os.Getwd()
	dir := ChdirTemp(t)

	current, _ := os.Getwd()
	if current != dir {
		t.Errorf("got %q, want %q", current, dir)
	}

	// Verify it's a temp directory
	if !filepath.IsAbs(dir) {
		t.Error("expected absolute path")
	}

	os.Chdir(original)
}

func TestAssertLen(t *testing.T) {
	AssertLen(t, []int{1, 2, 3}, 3)
	AssertLen(t, "hello", 5)
	AssertLen(t, map[string]int{"a": 1}, 1)
}

func TestAssertEmpty(t *testing.T) {
	AssertEmpty(t, []int{})
	AssertEmpty(t, "")
	AssertEmpty(t, map[string]int{})
}

func TestAssertNotEmpty(t *testing.T) {
	AssertNotEmpty(t, []int{1})
	AssertNotEmpty(t, "hello")
	AssertNotEmpty(t, map[string]int{"a": 1})
}
