// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDirectory_Override(t *testing.T) {
	customDir := "/custom/config/path"
	os.Setenv("GZH_CONFIG_DIR", customDir)
	defer os.Unsetenv("GZH_CONFIG_DIR")

	dir := GetConfigDirectory()
	if dir != customDir {
		t.Errorf("expected %s, got %s", customDir, dir)
	}
}

func TestGetConfigDirectory_Default(t *testing.T) {
	os.Unsetenv("GZH_CONFIG_DIR")

	dir := GetConfigDirectory()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir failed: %v", err)
	}
	expected := filepath.Join(homeDir, ".config", "gzh-manager")

	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestGetConfigDirectory_EmptyOverride(t *testing.T) {
	os.Setenv("GZH_CONFIG_DIR", "")
	defer os.Unsetenv("GZH_CONFIG_DIR")

	dir := GetConfigDirectory()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir failed: %v", err)
	}
	expected := filepath.Join(homeDir, ".config", "gzh-manager")

	if dir != expected {
		t.Errorf("empty GZH_CONFIG_DIR should fall back to default, expected %s, got %s", expected, dir)
	}
}

func TestEnsureConfigDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "gzh-test-config")
	os.Setenv("GZH_CONFIG_DIR", configDir)
	defer os.Unsetenv("GZH_CONFIG_DIR")

	if err := EnsureConfigDirectory(); err != nil {
		t.Fatalf("EnsureConfigDirectory failed: %v", err)
	}

	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected a directory")
	}
}

func TestEnsureConfigDirectory_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("GZH_CONFIG_DIR", tmpDir)
	defer os.Unsetenv("GZH_CONFIG_DIR")

	// Directory already exists — should be idempotent
	if err := EnsureConfigDirectory(); err != nil {
		t.Fatalf("EnsureConfigDirectory on existing dir failed: %v", err)
	}
}
