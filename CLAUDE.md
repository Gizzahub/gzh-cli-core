# CLAUDE.md

This file provides LLM-optimized guidance for Claude Code when working with this repository.

---

## Project Context

**Module**: `github.com/gizzahub/gzh-cli-core`
**Type**: Shared library for gzh-cli-* tools
**Go Version**: 1.24+

### Purpose

Core library providing common utilities across all gzh-cli tools:
- Logger, TestUtil, Errors, Config, CLI, Version

---

## Package Overview

| Package | Purpose | Key Types |
|---------|---------|-----------|
| `logger` | Structured logging | `Logger`, `SimpleLogger`, `Level` |
| `testutil` | Test helpers | `TempDir`, `Assert*`, `Capture` |
| `errors` | Error handling | Sentinel errors, `Wrap*`, validation |
| `config` | Config loading | `Loader`, env helpers |
| `cli` | Cobra helpers | `GlobalFlags`, `Output` |
| `version` | Version info | `Info`, `Get()` |

---

## Development Workflow

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Format code
go fmt ./...
```

---

## Important Rules

### Adding New Utilities

1. Keep packages focused - one concern per package
2. Write tests for all public functions
3. Use interfaces for testability
4. Document public APIs

### Dependencies

- Minimize external dependencies
- Only stdlib + yaml.v3 + cobra/pflag
- No CGO dependencies (cross-compile friendly)

### Commit Format

```
{type}({scope}): {description}

Model: claude-{model}
Co-Authored-By: Claude <noreply@anthropic.com>
```

---

## Project Structure

```
gzh-cli-core/
├── logger/          # Structured logging
│   ├── logger.go    # Logger interface & SimpleLogger
│   ├── global.go    # Default logger functions
│   └── logger_test.go
├── testutil/        # Test utilities
│   ├── testutil.go  # TempDir, Capture
│   ├── assert.go    # Assertions
│   ├── env.go       # SetEnv, Chdir
│   └── testutil_test.go
├── errors/          # Error handling
│   ├── errors.go    # Core error types & wrapping
│   ├── validation.go # Validation error helpers
│   └── errors_test.go
├── config/          # Configuration
│   ├── loader.go    # YAML config loader
│   ├── env.go       # Environment variable helpers
│   └── config_test.go
├── cli/             # CLI utilities
│   ├── flags.go     # Common flag sets
│   ├── root.go      # Root command helpers
│   ├── output.go    # Output formatting
│   └── cli_test.go
├── version/         # Version info
│   ├── version.go   # Version struct & Get()
│   └── version_test.go
├── go.mod
├── README.md
└── CLAUDE.md
```

---

## Usage in Other Projects

```go
// go.mod
require github.com/gizzahub/gzh-cli-core v0.1.0

// Import packages as needed
import (
    "github.com/gizzahub/gzh-cli-core/logger"
    "github.com/gizzahub/gzh-cli-core/errors"
    "github.com/gizzahub/gzh-cli-core/config"
)
```

---

**Last Updated**: 2024-12-05
