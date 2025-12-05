package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	l := New("test")
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	if l.name != "test" {
		t.Errorf("expected name 'test', got '%s'", l.name)
	}
	if l.level != LevelInfo {
		t.Errorf("expected default level Info, got %v", l.level)
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("Level(%d).String() = %s, want %s", tt.level, got, tt.want)
		}
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"unknown", LevelInfo}, // default
	}

	for _, tt := range tests {
		if got := ParseLevel(tt.input); got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestSimpleLogger_SetLevel(t *testing.T) {
	l := New("test")
	l.SetLevel(LevelDebug)
	if l.GetLevel() != LevelDebug {
		t.Errorf("expected LevelDebug, got %v", l.GetLevel())
	}
}

func TestSimpleLogger_Logging(t *testing.T) {
	var buf bytes.Buffer
	l := New("test")
	l.SetOutput(&buf)
	l.SetLevel(LevelDebug)

	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Error("expected DEBUG in output")
	}
	if !strings.Contains(output, "INFO") {
		t.Error("expected INFO in output")
	}
	if !strings.Contains(output, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("expected ERROR in output")
	}
	if !strings.Contains(output, "[test]") {
		t.Error("expected logger name in output")
	}
}

func TestSimpleLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := New("test")
	l.SetOutput(&buf)
	l.SetLevel(LevelWarn)

	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")

	output := buf.String()
	if strings.Contains(output, "DEBUG") {
		t.Error("DEBUG should be filtered out")
	}
	if strings.Contains(output, "INFO") {
		t.Error("INFO should be filtered out")
	}
	if !strings.Contains(output, "WARN") {
		t.Error("expected WARN in output")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("expected ERROR in output")
	}
}

func TestSimpleLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	l := New("test")
	l.SetOutput(&buf)

	contextLogger := l.WithContext("request_id", "abc123")
	contextLogger.Info("with context")

	output := buf.String()
	if !strings.Contains(output, "request_id=abc123") {
		t.Errorf("expected context in output, got: %s", output)
	}
}

func TestSimpleLogger_KeyValueArgs(t *testing.T) {
	var buf bytes.Buffer
	l := New("test")
	l.SetOutput(&buf)

	l.Info("message", "key1", "value1", "key2", 42)

	output := buf.String()
	if !strings.Contains(output, "key1=value1") {
		t.Errorf("expected key1=value1 in output, got: %s", output)
	}
	if !strings.Contains(output, "key2=42") {
		t.Errorf("expected key2=42 in output, got: %s", output)
	}
}

func TestNopLogger(t *testing.T) {
	l := NewNop()

	// Should not panic
	l.Debug("test")
	l.Info("test")
	l.Warn("test")
	l.Error("test")
	l.SetLevel(LevelDebug)
	l.SetOutput(nil)

	// WithContext should return same type
	ctx := l.WithContext("key", "value")
	if _, ok := ctx.(*NopLogger); !ok {
		t.Error("WithContext should return NopLogger")
	}
}

func TestGlobalLogger(t *testing.T) {
	var buf bytes.Buffer
	original := Default()

	l := New("global-test")
	l.SetOutput(&buf)
	SetDefault(l)

	Info("global info")
	if !strings.Contains(buf.String(), "global info") {
		t.Error("expected global logger to work")
	}

	// Restore original
	SetDefault(original)
}
