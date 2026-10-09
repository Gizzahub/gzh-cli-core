# gzh-cli-core

This file provides LLM-optimized guidance for Claude Code when working with this repository.

---

## Project Context

**Module**: `github.com/gizzahub/gzh-cli-core`
**Type**: Shared library for gzh-cli-* tools
**Consumer Go Version**: 1.26 (`go.mod` directive)
**Preferred Development / CI Toolchain**: 1.26.7 (main module directive / CI pin; not a consumer floor)

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

**Last Updated**: 2026-08-31
