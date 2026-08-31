package config

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

type testConfig struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Debug   bool   `yaml:"debug"`
	Timeout string `yaml:"timeout"`
}

func TestNewLoader(t *testing.T) {
	l := NewLoader("myapp")
	if l == nil {
		t.Fatal("expected non-nil loader")
	}
	if l.appName != "myapp" {
		t.Errorf("expected appName 'myapp', got '%s'", l.appName)
	}
	if len(l.paths) == 0 {
		t.Error("expected default paths")
	}
}

func TestLoader_WithPaths(t *testing.T) {
	l := NewLoader("myapp").WithPaths("custom.yaml", "other.yaml")
	if len(l.paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(l.paths))
	}
}

func TestLoader_AddPath(t *testing.T) {
	l := NewLoader("myapp")
	originalLen := len(l.paths)
	l.AddPath("extra.yaml")
	if len(l.paths) != originalLen+1 {
		t.Error("expected path to be added")
	}
	if l.paths[len(l.paths)-1] != "extra.yaml" {
		t.Error("expected extra.yaml at end")
	}
}

func TestLoader_PrependPath(t *testing.T) {
	l := NewLoader("myapp")
	l.PrependPath("first.yaml")
	if l.paths[0] != "first.yaml" {
		t.Error("expected first.yaml at beginning")
	}
}

func TestLoader_LoadFrom(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `name: test-app
port: 8080
debug: true
timeout: 30s`

	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	l := NewLoader("myapp")
	var cfg testConfig
	if err := l.LoadFrom(configPath, &cfg); err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if cfg.Name != "test-app" {
		t.Errorf("expected name 'test-app', got '%s'", cfg.Name)
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if !cfg.Debug {
		t.Error("expected debug true")
	}
}

func TestLoader_LoadFromErrors(t *testing.T) {
	l := NewLoader("myapp")

	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.yaml")
		err := l.LoadFrom(path, &testConfig{})
		if err == nil || !strings.Contains(err.Error(), "failed to read config file "+path) {
			t.Fatalf("LoadFrom() error = %v, want read error for %q", err, path)
		}
	})

	t.Run("malformed YAML", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "malformed.yaml")
		if err := os.WriteFile(path, []byte("name: ["), 0o644); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		err := l.LoadFrom(path, &testConfig{})
		if err == nil || !strings.Contains(err.Error(), "failed to parse config file "+path) {
			t.Fatalf("LoadFrom() error = %v, want parse error for %q", err, path)
		}
	})
}

func TestLoader_Load(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "myapp.yaml")

	content := `name: loaded-app`
	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	l := NewLoader("myapp").WithPaths(configPath)
	var cfg testConfig
	if err := l.Load(&cfg); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Name != "loaded-app" {
		t.Errorf("expected name 'loaded-app', got '%s'", cfg.Name)
	}
}

func TestLoader_Load_NotFound(t *testing.T) {
	l := NewLoader("nonexistent").WithPaths("nonexistent1.yaml", "nonexistent2.yaml")
	var cfg testConfig
	err := l.Load(&cfg)
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func TestLoader_LoadOrDefault(t *testing.T) {
	l := NewLoader("nonexistent").WithPaths("nonexistent.yaml")

	cfg := testConfig{
		Name: "default-name",
		Port: 3000,
	}

	if err := l.LoadOrDefault(&cfg); err != nil {
		t.Fatalf("LoadOrDefault failed: %v", err)
	}

	// Should retain defaults
	if cfg.Name != "default-name" {
		t.Errorf("expected default name, got '%s'", cfg.Name)
	}
	if cfg.Port != 3000 {
		t.Errorf("expected default port, got %d", cfg.Port)
	}
}

func TestLoader_LoadOrDefaultLoadsFirstExistingPath(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("name: configured\nport: 8080\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	l := NewLoader("myapp").WithPaths(filepath.Join(dir, "missing.yaml"), configPath)
	cfg := testConfig{Name: "default", Port: 3000}
	if err := l.LoadOrDefault(&cfg); err != nil {
		t.Fatalf("LoadOrDefault() error = %v", err)
	}
	if got, want := cfg.Name, "configured"; got != want {
		t.Errorf("Name = %q, want %q", got, want)
	}
	if got, want := cfg.Port, 8080; got != want {
		t.Errorf("Port = %d, want %d", got, want)
	}
}

func TestLoader_FindConfigFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "found.yaml")
	if err := os.WriteFile(configPath, []byte("test: value"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	l := NewLoader("app").WithPaths("notfound.yaml", configPath)
	path, found := l.FindConfigFile()
	if !found {
		t.Error("expected to find config file")
	}
	if path != configPath {
		t.Errorf("expected path '%s', got '%s'", configPath, path)
	}
}

func TestLoader_FindConfigFileNotFound(t *testing.T) {
	l := NewLoader("app").WithPaths(filepath.Join(t.TempDir(), "missing.yaml"))
	if path, found := l.FindConfigFile(); found || path != "" {
		t.Errorf("FindConfigFile() = (%q, %t), want (\"\", false)", path, found)
	}
}

func TestSave(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "subdir", "config.yaml")

	cfg := testConfig{
		Name:  "saved-app",
		Port:  9000,
		Debug: true,
	}

	if err := Save(configPath, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and can be loaded
	l := NewLoader("app")
	var loaded testConfig
	if err := l.LoadFrom(configPath, &loaded); err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}
	if loaded.Name != "saved-app" {
		t.Errorf("expected name 'saved-app', got '%s'", loaded.Name)
	}
}

func TestDefaultPaths(t *testing.T) {
	paths := DefaultPaths("testapp")
	if len(paths) == 0 {
		t.Error("expected non-empty paths")
	}

	// Should contain app-specific paths
	foundYaml := slices.Contains(paths, "testapp.yaml")
	if !foundYaml {
		t.Error("expected 'testapp.yaml' in paths")
	}
}

// Environment variable tests

func TestGetEnv(t *testing.T) {
	os.Setenv("GZH_TEST_VAR", "test_value")
	defer os.Unsetenv("GZH_TEST_VAR")

	v := GetEnv("TEST_VAR")
	if v != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", v)
	}
}

func TestGetEnvOr(t *testing.T) {
	v := GetEnvOr("NONEXISTENT_VAR", "default")
	if v != "default" {
		t.Errorf("expected 'default', got '%s'", v)
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"1", true},
		{"yes", true},
		{"on", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"", false},
	}

	for _, tt := range tests {
		os.Setenv("GZH_BOOL_TEST", tt.value)
		got := GetEnvBool("BOOL_TEST")
		if got != tt.expected {
			t.Errorf("GetEnvBool(%q) = %v, want %v", tt.value, got, tt.expected)
		}
	}
	os.Unsetenv("GZH_BOOL_TEST")
}

func TestGetEnvBoolOr(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue bool
		want         bool
	}{
		{name: "unset uses default", defaultValue: true, want: true},
		{name: "recognized true", value: "yes", want: true},
		{name: "recognized false", value: "false", defaultValue: true, want: false},
		{name: "invalid is false", value: "not-a-bool", defaultValue: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GZH_BOOL_OR_TEST", tt.value)
			if got := GetEnvBoolOr("BOOL_OR_TEST", tt.defaultValue); got != tt.want {
				t.Errorf("GetEnvBoolOr() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("GZH_INT_TEST", "42")
	defer os.Unsetenv("GZH_INT_TEST")

	v, ok := GetEnvInt("INT_TEST")
	if !ok {
		t.Error("expected ok to be true")
	}
	if v != 42 {
		t.Errorf("expected 42, got %d", v)
	}
}

func TestGetEnvIntOr(t *testing.T) {
	v := GetEnvIntOr("NONEXISTENT_INT", 100)
	if v != 100 {
		t.Errorf("expected 100, got %d", v)
	}
}

func TestGetEnvDuration(t *testing.T) {
	os.Setenv("GZH_DUR_TEST", "5m30s")
	defer os.Unsetenv("GZH_DUR_TEST")

	d, ok := GetEnvDuration("DUR_TEST")
	if !ok {
		t.Error("expected ok to be true")
	}
	expected := 5*time.Minute + 30*time.Second
	if d != expected {
		t.Errorf("expected %v, got %v", expected, d)
	}
}

func TestGetEnvNumericAndDurationInvalidValues(t *testing.T) {
	t.Setenv("GZH_INT_INVALID", "not-an-int")
	if value, ok := GetEnvInt("INT_INVALID"); ok || value != 0 {
		t.Errorf("GetEnvInt() = (%d, %t), want (0, false)", value, ok)
	}

	t.Setenv("GZH_DURATION_INVALID", "not-a-duration")
	if value, ok := GetEnvDuration("DURATION_INVALID"); ok || value != 0 {
		t.Errorf("GetEnvDuration() = (%s, %t), want (0s, false)", value, ok)
	}
	if got, want := GetEnvDurationOr("DURATION_INVALID", time.Minute), time.Minute; got != want {
		t.Errorf("GetEnvDurationOr() = %s, want %s", got, want)
	}
}

func TestGetEnvList(t *testing.T) {
	os.Setenv("GZH_LIST_TEST", "a, b, c")
	defer os.Unsetenv("GZH_LIST_TEST")

	list := GetEnvList("LIST_TEST")
	if len(list) != 3 {
		t.Errorf("expected 3 items, got %d", len(list))
	}
	if list[0] != "a" || list[1] != "b" || list[2] != "c" {
		t.Errorf("unexpected list: %v", list)
	}
}

func TestGetEnvListTrimsEmptyEntriesWithNoPrefix(t *testing.T) {
	t.Setenv("LIST_WITHOUT_PREFIX", " first, , second ,, ")
	if got, want := GetEnvList("LIST_WITHOUT_PREFIX", ""), []string{"first", "second"}; !slices.Equal(got, want) {
		t.Errorf("GetEnvList() = %v, want %v", got, want)
	}
}

func TestGetEnv_CustomPrefix(t *testing.T) {
	os.Setenv("MYAPP_CUSTOM", "custom_value")
	defer os.Unsetenv("MYAPP_CUSTOM")

	v := GetEnv("CUSTOM", "MYAPP")
	if v != "custom_value" {
		t.Errorf("expected 'custom_value', got '%s'", v)
	}
}

func TestLookupEnv(t *testing.T) {
	os.Setenv("GZH_LOOKUP_TEST", "")
	defer os.Unsetenv("GZH_LOOKUP_TEST")

	// Variable is set but empty
	v, ok := LookupEnv("LOOKUP_TEST")
	if !ok {
		t.Error("expected ok to be true for set but empty var")
	}
	if v != "" {
		t.Errorf("expected empty string, got '%s'", v)
	}

	// Variable not set
	_, ok = LookupEnv("DEFINITELY_NOT_SET_VAR")
	if ok {
		t.Error("expected ok to be false for unset var")
	}
}

func TestMustGetEnv(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		t.Setenv("GZH_REQUIRED", "value")
		if got, want := MustGetEnv("REQUIRED"), "value"; got != want {
			t.Errorf("MustGetEnv() = %q, want %q", got, want)
		}
	})

	t.Run("unset panics with full key", func(t *testing.T) {
		t.Setenv("GZH_REQUIRED", "")
		defer func() {
			if recovered := recover(); recovered != "required environment variable not set: GZH_REQUIRED" {
				t.Errorf("panic = %v, want required-key message", recovered)
			}
		}()
		MustGetEnv("REQUIRED")
	})
}
