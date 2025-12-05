package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("expected non-empty version")
	}
	if info.GoVersion == "" {
		t.Error("expected non-empty go version")
	}
	if info.Platform == "" {
		t.Error("expected non-empty platform")
	}

	// Platform should contain GOOS/GOARCH
	expectedPlatform := runtime.GOOS + "/" + runtime.GOARCH
	if info.Platform != expectedPlatform {
		t.Errorf("expected platform '%s', got '%s'", expectedPlatform, info.Platform)
	}
}

func TestInfo_String(t *testing.T) {
	info := Info{
		Version:   "1.2.3",
		GitCommit: "abc1234567890",
		BuildDate: "2024-06-15",
		GoVersion: "go1.21.0",
		Platform:  "linux/amd64",
	}

	s := info.String()

	checks := []string{
		"1.2.3",
		"abc1234567890",
		"2024-06-15",
		"go1.21.0",
		"linux/amd64",
	}

	for _, check := range checks {
		if !strings.Contains(s, check) {
			t.Errorf("expected '%s' in output:\n%s", check, s)
		}
	}
}

func TestInfo_String_Unknown(t *testing.T) {
	info := Info{
		Version:   "dev",
		GitCommit: "unknown",
		BuildDate: "unknown",
		GoVersion: "go1.21.0",
		Platform:  "darwin/arm64",
	}

	s := info.String()

	// Should not include "unknown" values in formatted output
	if strings.Contains(s, "Git Commit: unknown") {
		t.Error("should not display unknown git commit")
	}
	if strings.Contains(s, "Build Date: unknown") {
		t.Error("should not display unknown build date")
	}
}

func TestInfo_Short(t *testing.T) {
	info := Info{Version: "2.0.0"}
	if info.Short() != "2.0.0" {
		t.Errorf("expected '2.0.0', got '%s'", info.Short())
	}
}

func TestInfo_Full(t *testing.T) {
	tests := []struct {
		name      string
		info      Info
		expected  string
	}{
		{
			name:     "with commit",
			info:     Info{Version: "1.0.0", GitCommit: "abc1234567890"},
			expected: "1.0.0-abc1234",
		},
		{
			name:     "unknown commit",
			info:     Info{Version: "1.0.0", GitCommit: "unknown"},
			expected: "1.0.0",
		},
		{
			name:     "empty commit",
			info:     Info{Version: "1.0.0", GitCommit: ""},
			expected: "1.0.0",
		},
		{
			name:     "short commit",
			info:     Info{Version: "1.0.0", GitCommit: "abc"},
			expected: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.info.Full()
			if got != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

func TestLdFlags(t *testing.T) {
	flags := LdFlags("main", "1.0.0", "abc123", "2024-01-01")

	if !strings.Contains(flags, "-X main.Version=1.0.0") {
		t.Error("expected Version ldflags")
	}
	if !strings.Contains(flags, "-X main.GitCommit=abc123") {
		t.Error("expected GitCommit ldflags")
	}
	if !strings.Contains(flags, "-X main.BuildDate=2024-01-01") {
		t.Error("expected BuildDate ldflags")
	}
}
