package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func unsetEnvForTest(t *testing.T, key string) {
	t.Helper()
	previous, wasSet := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv(%q) error = %v", key, err)
	}
	t.Cleanup(func() {
		var err error
		if wasSet {
			err = os.Setenv(key, previous)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Errorf("restore %q: %v", key, err)
		}
	})
}

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

func TestCaptureSplitsOutputAndRestoresStreams(t *testing.T) {
	stdout, stderr := os.Stdout, os.Stderr
	output := Capture(func() {
		fmt.Fprintln(os.Stdout, "stdout message")
		fmt.Fprintln(os.Stderr, "stderr message")
	})

	if got, want := output.Stdout, "stdout message\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if got, want := output.Stderr, "stderr message\n"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
	if os.Stdout != stdout || os.Stderr != stderr {
		t.Error("Capture() did not restore stdout and stderr")
	}
}

func TestCaptureStdoutRestoresStream(t *testing.T) {
	stdout := os.Stdout
	output := CaptureStdout(func() {
		fmt.Println("stdout only")
	})
	if got, want := output, "stdout only\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if os.Stdout != stdout {
		t.Error("CaptureStdout() did not restore stdout")
	}
}

func TestCaptureStderrRestoresStream(t *testing.T) {
	stderr := os.Stderr
	output := CaptureStderr(func() {
		fmt.Fprintln(os.Stderr, "stderr only")
	})
	if got, want := output, "stderr only\n"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
	if os.Stderr != stderr {
		t.Error("CaptureStderr() did not restore stderr")
	}
}

func TestAssertNoError(t *testing.T) {
	// Should not fail
	AssertNoError(t, nil)
}

func TestAssertError(t *testing.T) {
	// Should not fail
	AssertError(t, errors.New("some error"))
}

func TestAdditionalAssertionSuccessCases(t *testing.T) {
	AssertErrorContains(t, errors.New("connection refused"), "refused")
	AssertNotEqual(t, "left", "right")
	AssertTrue(t, true, "expected true")
	AssertFalse(t, false, "expected false")

	var typedNil *string
	AssertNil(t, nil)
	AssertNil(t, typedNil)
	AssertNotNil(t, &typedNil)
	AssertNotNil(t, []string{})
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

func TestSetEnvRestoresExistingAndUnsetState(t *testing.T) {
	t.Run("existing value", func(t *testing.T) {
		const key = "GZH_TESTUTIL_SET_EXISTING"
		t.Setenv(key, "original")
		t.Run("set", func(t *testing.T) {
			SetEnv(t, key, "test_value")
			if got, want := os.Getenv(key), "test_value"; got != want {
				t.Errorf("value = %q, want %q", got, want)
			}
		})
		if got, want := os.Getenv(key), "original"; got != want {
			t.Errorf("restored value = %q, want %q", got, want)
		}
	})

	t.Run("unset value", func(t *testing.T) {
		const key = "GZH_TESTUTIL_SET_UNSET"
		unsetEnvForTest(t, key)
		t.Run("set", func(t *testing.T) {
			SetEnv(t, key, "test_value")
			if got, want := os.Getenv(key), "test_value"; got != want {
				t.Errorf("value = %q, want %q", got, want)
			}
		})
		if _, exists := os.LookupEnv(key); exists {
			t.Errorf("%s remained set after cleanup", key)
		}
	})
}

func TestUnsetEnvRestoresExistingAndUnsetState(t *testing.T) {
	t.Run("existing value", func(t *testing.T) {
		const key = "GZH_TESTUTIL_UNSET_EXISTING"
		t.Setenv(key, "original")
		t.Run("unset", func(t *testing.T) {
			UnsetEnv(t, key)
			if _, exists := os.LookupEnv(key); exists {
				t.Errorf("%s remained set", key)
			}
		})
		if got, want := os.Getenv(key), "original"; got != want {
			t.Errorf("restored value = %q, want %q", got, want)
		}
	})

	t.Run("unset value", func(t *testing.T) {
		const key = "GZH_TESTUTIL_UNSET_UNSET"
		unsetEnvForTest(t, key)
		t.Run("unset", func(t *testing.T) {
			UnsetEnv(t, key)
		})
		if _, exists := os.LookupEnv(key); exists {
			t.Errorf("%s became set after cleanup", key)
		}
	})
}

func TestChdirAndChdirTempRestoreWorkingDirectory(t *testing.T) {
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}

	t.Run("Chdir", func(t *testing.T) {
		dir := TempDir(t)
		Chdir(t, dir)
		if current, err := os.Getwd(); err != nil || current != dir {
			t.Errorf("Getwd() = (%q, %v), want (%q, nil)", current, err, dir)
		}
	})
	if current, err := os.Getwd(); err != nil || current != original {
		t.Errorf("Getwd() after Chdir cleanup = (%q, %v), want (%q, nil)", current, err, original)
	}

	t.Run("ChdirTemp", func(t *testing.T) {
		dir := ChdirTemp(t)
		if current, err := os.Getwd(); err != nil || current != dir {
			t.Errorf("Getwd() = (%q, %v), want (%q, nil)", current, err, dir)
		}
		if !filepath.IsAbs(dir) {
			t.Error("ChdirTemp() directory is not absolute")
		}
	})
	if current, err := os.Getwd(); err != nil || current != original {
		t.Errorf("Getwd() after ChdirTemp cleanup = (%q, %v), want (%q, nil)", current, err, original)
	}
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
