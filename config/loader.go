// Package config provides configuration loading utilities for gzh-cli tools.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Loader provides configuration file loading utilities.
type Loader struct {
	appName string
	paths   []string
}

// NewLoader creates a new configuration loader with the given app name.
func NewLoader(appName string) *Loader {
	return &Loader{
		appName: appName,
		paths:   DefaultPaths(appName),
	}
}

// WithPaths sets custom search paths.
func (l *Loader) WithPaths(paths ...string) *Loader {
	l.paths = paths
	return l
}

// AddPath adds a path to the search list.
func (l *Loader) AddPath(path string) *Loader {
	l.paths = append(l.paths, path)
	return l
}

// PrependPath adds a path to the beginning of the search list.
func (l *Loader) PrependPath(path string) *Loader {
	l.paths = append([]string{path}, l.paths...)
	return l
}

// Paths returns the current search paths.
func (l *Loader) Paths() []string {
	return l.paths
}

// Load loads configuration from the first existing file in the search paths.
// The dst must be a pointer to a struct.
func (l *Loader) Load(dst interface{}) error {
	for _, path := range l.paths {
		if _, err := os.Stat(path); err == nil {
			return l.LoadFrom(path, dst)
		}
	}
	return fmt.Errorf("no config file found in paths: %v", l.paths)
}

// LoadFrom loads configuration from a specific file path.
func (l *Loader) LoadFrom(path string, dst interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return nil
}

// LoadOrDefault loads configuration, returning nil error if no file found.
// Caller should initialize dst with default values before calling.
func (l *Loader) LoadOrDefault(dst interface{}) error {
	for _, path := range l.paths {
		if _, err := os.Stat(path); err == nil {
			return l.LoadFrom(path, dst)
		}
	}
	// No config file found, dst retains its default values
	return nil
}

// FindConfigFile returns the first existing config file path.
func (l *Loader) FindConfigFile() (string, bool) {
	for _, path := range l.paths {
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

// DefaultPaths returns default configuration search paths for the given app.
func DefaultPaths(appName string) []string {
	paths := []string{
		appName + ".yaml",
		appName + ".yml",
		"." + appName + ".yaml",
		"." + appName + ".yml",
	}

	// Add XDG config path
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		paths = append(paths,
			filepath.Join(configDir, appName, "config.yaml"),
			filepath.Join(configDir, appName, "config.yml"),
		)
	}

	// Add home directory paths
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths,
			filepath.Join(home, ".config", appName, "config.yaml"),
			filepath.Join(home, ".config", appName, "config.yml"),
			filepath.Join(home, "."+appName+".yaml"),
			filepath.Join(home, "."+appName+".yml"),
		)
	}

	return paths
}

// Save saves configuration to the given path.
func Save(path string, cfg interface{}) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
