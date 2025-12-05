// Package version provides version information utilities for gzh-cli tools.
package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time via ldflags.
var (
	// Version is the semantic version (e.g., "1.0.0")
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the build timestamp
	BuildDate = "unknown"
)

// Info holds version information.
type Info struct {
	Version   string `json:"version" yaml:"version"`
	GitCommit string `json:"git_commit" yaml:"git_commit"`
	BuildDate string `json:"build_date" yaml:"build_date"`
	GoVersion string `json:"go_version" yaml:"go_version"`
	Platform  string `json:"platform" yaml:"platform"`
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string.
func (i Info) String() string {
	s := fmt.Sprintf("Version:    %s", i.Version)
	if i.GitCommit != "" && i.GitCommit != "unknown" {
		s += fmt.Sprintf("\nGit Commit: %s", i.GitCommit)
	}
	if i.BuildDate != "" && i.BuildDate != "unknown" {
		s += fmt.Sprintf("\nBuild Date: %s", i.BuildDate)
	}
	s += fmt.Sprintf("\nGo Version: %s", i.GoVersion)
	s += fmt.Sprintf("\nPlatform:   %s", i.Platform)
	return s
}

// Short returns just the version string.
func (i Info) Short() string {
	return i.Version
}

// Full returns version with git commit (e.g., "1.0.0-abc1234")
func (i Info) Full() string {
	if i.GitCommit != "" && i.GitCommit != "unknown" && len(i.GitCommit) >= 7 {
		return fmt.Sprintf("%s-%s", i.Version, i.GitCommit[:7])
	}
	return i.Version
}

// LdFlags returns the ldflags string for building with version info.
// Usage: go build -ldflags "$(version.LdFlags(pkg, ver, commit, date))"
func LdFlags(pkg, version, gitCommit, buildDate string) string {
	return fmt.Sprintf("-X %s.Version=%s -X %s.GitCommit=%s -X %s.BuildDate=%s",
		pkg, version, pkg, gitCommit, pkg, buildDate)
}
