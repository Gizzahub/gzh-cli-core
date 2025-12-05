package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootConfig holds configuration for creating a root command.
type RootConfig struct {
	// Name is the binary name (e.g., "gz-git")
	Name string
	// Short is a short description
	Short string
	// Long is a long description
	Long string
	// Version string
	Version string
	// VersionTemplate customizes version output
	VersionTemplate string
}

// NewRootCmd creates a new root command with standard configuration.
func NewRootCmd(cfg RootConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   cfg.Name,
		Short: cfg.Short,
		Long:  cfg.Long,
	}

	if cfg.Version != "" {
		cmd.Version = cfg.Version
		if cfg.VersionTemplate != "" {
			cmd.SetVersionTemplate(cfg.VersionTemplate)
		} else {
			cmd.SetVersionTemplate(fmt.Sprintf("%s version {{.Version}}\n", cfg.Name))
		}
	}

	return cmd
}

// Execute runs the root command and handles errors.
func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// ExecuteWithCode runs the root command and returns the exit code.
func ExecuteWithCode(cmd *cobra.Command) int {
	if err := cmd.Execute(); err != nil {
		return 1
	}
	return 0
}

// AddVersionCmd adds a version subcommand with extended info.
func AddVersionCmd(root *cobra.Command, info VersionInfo) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(info.String())
		},
	}
	root.AddCommand(versionCmd)
}

// VersionInfo holds extended version information.
type VersionInfo struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
	Platform  string
}

// String returns a formatted version string.
func (v VersionInfo) String() string {
	s := fmt.Sprintf("Version:    %s", v.Version)
	if v.GitCommit != "" && v.GitCommit != "unknown" {
		s += fmt.Sprintf("\nGit Commit: %s", v.GitCommit)
	}
	if v.BuildDate != "" && v.BuildDate != "unknown" {
		s += fmt.Sprintf("\nBuild Date: %s", v.BuildDate)
	}
	if v.GoVersion != "" {
		s += fmt.Sprintf("\nGo Version: %s", v.GoVersion)
	}
	if v.Platform != "" {
		s += fmt.Sprintf("\nPlatform:   %s", v.Platform)
	}
	return s
}

// Short returns just the version string.
func (v VersionInfo) Short() string {
	return v.Version
}
