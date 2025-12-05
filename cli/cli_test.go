package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGlobalFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &GlobalFlags{}
	AddGlobalFlags(cmd, flags)

	// Test that flags are registered
	if cmd.PersistentFlags().Lookup("verbose") == nil {
		t.Error("expected verbose flag")
	}
	if cmd.PersistentFlags().Lookup("quiet") == nil {
		t.Error("expected quiet flag")
	}
	if cmd.PersistentFlags().Lookup("debug") == nil {
		t.Error("expected debug flag")
	}
	if cmd.PersistentFlags().Lookup("no-color") == nil {
		t.Error("expected no-color flag")
	}
	if cmd.PersistentFlags().Lookup("config") == nil {
		t.Error("expected config flag")
	}
}

func TestOutputFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &OutputFlags{}
	AddOutputFlags(cmd, flags)

	if cmd.Flags().Lookup("format") == nil {
		t.Error("expected format flag")
	}
	if cmd.Flags().Lookup("output") == nil {
		t.Error("expected output flag")
	}
}

func TestDryRunFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	flags := &DryRunFlags{}
	AddDryRunFlags(cmd, flags)

	if cmd.Flags().Lookup("dry-run") == nil {
		t.Error("expected dry-run flag")
	}
	if cmd.Flags().Lookup("force") == nil {
		t.Error("expected force flag")
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd(RootConfig{
		Name:    "test-app",
		Short:   "Test application",
		Long:    "A longer description",
		Version: "1.0.0",
	})

	if cmd.Use != "test-app" {
		t.Errorf("expected Use 'test-app', got '%s'", cmd.Use)
	}
	if cmd.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got '%s'", cmd.Version)
	}
}

func TestVersionInfo_String(t *testing.T) {
	info := VersionInfo{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2024-01-01",
		GoVersion: "1.21.0",
		Platform:  "linux/amd64",
	}

	s := info.String()
	if !strings.Contains(s, "1.0.0") {
		t.Error("expected version in output")
	}
	if !strings.Contains(s, "abc123") {
		t.Error("expected git commit in output")
	}
	if !strings.Contains(s, "2024-01-01") {
		t.Error("expected build date in output")
	}
}

func TestVersionInfo_Short(t *testing.T) {
	info := VersionInfo{Version: "2.0.0"}
	if info.Short() != "2.0.0" {
		t.Errorf("expected '2.0.0', got '%s'", info.Short())
	}
}

func TestOutput_Success(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf)
	out.Success("operation completed")

	if !strings.Contains(buf.String(), "✓") {
		t.Error("expected checkmark in success output")
	}
	if !strings.Contains(buf.String(), "operation completed") {
		t.Error("expected message in output")
	}
}

func TestOutput_Error(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf)
	out.Error("operation failed")

	if !strings.Contains(buf.String(), "✗") {
		t.Error("expected X mark in error output")
	}
}

func TestOutput_Warning(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf)
	out.Warning("be careful")

	if !strings.Contains(buf.String(), "⚠") {
		t.Error("expected warning symbol in output")
	}
}

func TestOutput_JSON(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("json")

	data := map[string]string{"key": "value"}
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if !strings.Contains(buf.String(), `"key"`) {
		t.Error("expected JSON key in output")
	}
}

func TestOutput_YAML(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf).SetFormat("yaml")

	data := map[string]string{"key": "value"}
	if err := out.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	if !strings.Contains(buf.String(), "key:") {
		t.Error("expected YAML key in output")
	}
}

func TestOutput_DryRun(t *testing.T) {
	var buf bytes.Buffer
	out := NewOutput().SetWriter(&buf)
	out.DryRun()

	if !strings.Contains(buf.String(), "[DRY-RUN]") {
		t.Error("expected DRY-RUN marker in output")
	}
}
