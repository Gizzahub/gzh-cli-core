package testutil

import (
	"os"
	"testing"
)

// SetEnv sets an environment variable and restores it when the test finishes.
func SetEnv(t *testing.T, key, value string) {
	t.Helper()
	old, existed := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set env %s: %v", key, err)
	}
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	})
}

// UnsetEnv unsets an environment variable and restores it when the test finishes.
func UnsetEnv(t *testing.T, key string) {
	t.Helper()
	old, existed := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset env %s: %v", key, err)
	}
	t.Cleanup(func() {
		if existed {
			os.Setenv(key, old)
		}
	})
}

// Chdir changes the working directory and restores it when the test finishes.
func Chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change directory to %s: %v", dir, err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(old); err != nil {
			t.Logf("warning: failed to restore directory to %s: %v", old, err)
		}
	})
}

// ChdirTemp creates a temp directory, changes to it, and restores when done.
func ChdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	Chdir(t, dir)
	return dir
}
