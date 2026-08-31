package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/gizzahub/gzh-cli-core/testutil"
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

func TestConfirmFlags_ParseLongAndShort(t *testing.T) {
	for _, arg := range []string{"--yes", "-y"} {
		t.Run(arg, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			flags := &ConfirmFlags{}
			AddConfirmFlags(cmd, flags)
			cmd.SetArgs([]string{arg})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !flags.Yes {
				t.Fatalf("%s did not set ConfirmFlags.Yes", arg)
			}
		})
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

func TestNewRootCmdVersionBehavior(t *testing.T) {
	t.Run("custom template", func(t *testing.T) {
		cmd := NewRootCmd(RootConfig{
			Name:            "test-app",
			Version:         "1.2.3",
			VersionTemplate: "release {{.Version}}\n",
		})
		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetArgs([]string{"--version"})

		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
		if got, want := output.String(), "release 1.2.3\n"; got != want {
			t.Errorf("version output = %q, want %q", got, want)
		}
	})

	t.Run("no version", func(t *testing.T) {
		cmd := NewRootCmd(RootConfig{Name: "test-app"})
		if cmd.Version != "" {
			t.Errorf("Version = %q, want empty", cmd.Version)
		}
	})
}

func TestExecuteWithCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cmd := &cobra.Command{RunE: func(*cobra.Command, []string) error { return nil }}
		if got := ExecuteWithCode(cmd); got != 0 {
			t.Errorf("ExecuteWithCode() = %d, want 0", got)
		}
	})

	t.Run("command error", func(t *testing.T) {
		cmd := &cobra.Command{RunE: func(*cobra.Command, []string) error { return errors.New("test failure") }}
		cmd.SetErr(&bytes.Buffer{})
		if got := ExecuteWithCode(cmd); got != 1 {
			t.Errorf("ExecuteWithCode() = %d, want 1", got)
		}
	})
}

func TestAddVersionCmd(t *testing.T) {
	root := NewRootCmd(RootConfig{Name: "test-app"})
	info := VersionInfo{
		Version:   "1.2.3",
		GitCommit: "abc123",
		BuildDate: "2026-08-31",
		GoVersion: "go1.26.7",
		Platform:  "linux/amd64",
	}
	AddVersionCmd(root, info)

	versionCmd, _, err := root.Find([]string{"version"})
	if err != nil {
		t.Fatalf("Find(version) error = %v", err)
	}
	if got, want := versionCmd.Short, "Print version information"; got != want {
		t.Errorf("version command Short = %q, want %q", got, want)
	}

	root.SetArgs([]string{"version"})
	var executeErr error
	output := testutil.CaptureStdout(func() {
		executeErr = root.Execute()
	})
	if executeErr != nil {
		t.Fatalf("Execute() error = %v", executeErr)
	}
	if got, want := output, info.String()+"\n"; got != want {
		t.Errorf("version output = %q, want %q", got, want)
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

func TestOutput_PrintTextFallback(t *testing.T) {
	for _, format := range []string{"text", "TEXT", "unknown"} {
		t.Run(format, func(t *testing.T) {
			var output bytes.Buffer
			out := NewOutput().SetWriter(&output).SetFormat(format)

			if err := out.Print("message"); err != nil {
				t.Fatalf("Print() error = %v", err)
			}
			if got, want := output.String(), "message\n"; got != want {
				t.Errorf("Print() output = %q, want %q", got, want)
			}
		})
	}
}

func TestOutput_InfoAndLine(t *testing.T) {
	var output bytes.Buffer
	out := NewOutput().SetWriter(&output)

	out.Info("ready: %s", "yes")
	out.Line("next")

	if got, want := output.String(), "ℹ ready: yes\nnext\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestPackageOutputHelpers(t *testing.T) {
	original := defaultOutput
	var output bytes.Buffer
	defaultOutput = NewOutput().SetWriter(&output)
	t.Cleanup(func() {
		defaultOutput = original
	})

	Success("done: %d", 1)
	Error("failed")
	Warning("careful")
	Info("ready")
	DryRun()

	const want = "✓ done: 1\n✗ failed\n⚠ careful\nℹ ready\n[DRY-RUN] No changes will be made\n"
	if got := output.String(); got != want {
		t.Errorf("package output = %q, want %q", got, want)
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
