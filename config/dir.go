// Copyright (c) 2025 Archmagece
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
)

// GetConfigDirectory returns the base configuration directory for all gzh-cli tools.
// It checks the GZH_CONFIG_DIR environment variable first (via GetEnv which prepends
// the DefaultEnvPrefix "GZH"), then falls back to ~/.config/gzh-manager.
func GetConfigDirectory() string {
	if dir := GetEnv("CONFIG_DIR"); dir != "" {
		return dir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(homeDir, ".config", "gzh-manager")
}

// EnsureConfigDirectory creates the configuration directory if it doesn't exist.
// The directory is created with 0o750 permissions.
func EnsureConfigDirectory() error {
	return os.MkdirAll(GetConfigDirectory(), 0o750)
}
