package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Output handles formatted output.
type Output struct {
	writer io.Writer
	format string
}

// NewOutput creates a new Output with default stdout writer.
func NewOutput() *Output {
	return &Output{
		writer: os.Stdout,
		format: "text",
	}
}

// SetWriter sets the output writer.
func (o *Output) SetWriter(w io.Writer) *Output {
	o.writer = w
	return o
}

// SetFormat sets the output format.
func (o *Output) SetFormat(format string) *Output {
	o.format = strings.ToLower(format)
	return o
}

// Print prints data in the configured format.
func (o *Output) Print(data interface{}) error {
	switch o.format {
	case "json":
		return o.printJSON(data)
	case "yaml", "yml":
		return o.printYAML(data)
	case "llm":
		return o.printLLM(data)
	default:
		return o.printText(data)
	}
}

// printLLM prints data in LLM-friendly compact format.
func (o *Output) printLLM(data interface{}) error {
	formatter := &llmFormatter{}
	output := formatter.format(data, 0)
	if output == "" {
		return nil
	}
	_, err := fmt.Fprint(o.writer, output)
	return err
}

func (o *Output) printJSON(data interface{}) error {
	enc := json.NewEncoder(o.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (o *Output) printYAML(data interface{}) error {
	enc := yaml.NewEncoder(o.writer)
	enc.SetIndent(2)
	return enc.Encode(data)
}

func (o *Output) printText(data interface{}) error {
	_, err := fmt.Fprintln(o.writer, data)
	return err
}

// Success prints a success message with checkmark.
func (o *Output) Success(msg string, args ...interface{}) {
	fmt.Fprintf(o.writer, "✓ "+msg+"\n", args...)
}

// Error prints an error message with X mark.
func (o *Output) Error(msg string, args ...interface{}) {
	fmt.Fprintf(o.writer, "✗ "+msg+"\n", args...)
}

// Warning prints a warning message.
func (o *Output) Warning(msg string, args ...interface{}) {
	fmt.Fprintf(o.writer, "⚠ "+msg+"\n", args...)
}

// Info prints an info message.
func (o *Output) Info(msg string, args ...interface{}) {
	fmt.Fprintf(o.writer, "ℹ "+msg+"\n", args...)
}

// Line prints a plain message.
func (o *Output) Line(msg string, args ...interface{}) {
	fmt.Fprintf(o.writer, msg+"\n", args...)
}

// DryRun prints a dry-run notice.
func (o *Output) DryRun() {
	fmt.Fprintln(o.writer, "[DRY-RUN] No changes will be made")
}

// Package-level convenience functions

var defaultOutput = NewOutput()

// Success prints a success message.
func Success(msg string, args ...interface{}) {
	defaultOutput.Success(msg, args...)
}

// Error prints an error message.
func Error(msg string, args ...interface{}) {
	defaultOutput.Error(msg, args...)
}

// Warning prints a warning message.
func Warning(msg string, args ...interface{}) {
	defaultOutput.Warning(msg, args...)
}

// Info prints an info message.
func Info(msg string, args ...interface{}) {
	defaultOutput.Info(msg, args...)
}

// DryRun prints a dry-run notice.
func DryRun() {
	defaultOutput.DryRun()
}
