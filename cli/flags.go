// Package cli provides CLI utilities for gzh-cli tools.
package cli

import (
	"github.com/spf13/cobra"
)

// GlobalFlags holds common flags used across all gzh-cli tools.
type GlobalFlags struct {
	Verbose bool
	Quiet   bool
	Debug   bool
	NoColor bool
	Config  string
}

// AddGlobalFlags adds common global flags to a command.
func AddGlobalFlags(cmd *cobra.Command, flags *GlobalFlags) {
	cmd.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "Suppress non-essential output")
	cmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug mode")
	cmd.PersistentFlags().BoolVar(&flags.NoColor, "no-color", false, "Disable colored output")
	cmd.PersistentFlags().StringVarP(&flags.Config, "config", "c", "", "Config file path")

	// Mark verbose and quiet as mutually exclusive
	cmd.MarkFlagsMutuallyExclusive("verbose", "quiet")
}

// OutputFlags holds flags for output formatting.
type OutputFlags struct {
	Format string // json, yaml, table, text
	Output string // output file path (empty for stdout)
}

// AddOutputFlags adds output formatting flags to a command.
func AddOutputFlags(cmd *cobra.Command, flags *OutputFlags) {
	cmd.Flags().StringVarP(&flags.Format, "format", "f", "text", "Output format (json, yaml, table, text)")
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "", "Output file (default: stdout)")
}

// DryRunFlags holds flags for dry-run mode.
type DryRunFlags struct {
	DryRun bool
	Force  bool
}

// AddDryRunFlags adds dry-run related flags to a command.
func AddDryRunFlags(cmd *cobra.Command, flags *DryRunFlags) {
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "Show what would be done without making changes")
	cmd.Flags().BoolVar(&flags.Force, "force", false, "Force operation without confirmation")
}

// ConfirmFlags holds flags for confirmation prompts.
type ConfirmFlags struct {
	Yes bool
}

// AddConfirmFlags adds confirmation flags to a command.
func AddConfirmFlags(cmd *cobra.Command, flags *ConfirmFlags) {
	cmd.Flags().BoolVarP(&flags.Yes, "yes", "y", false, "Assume yes to all prompts")
}
