# gzh-cli-core

## Development workflow — user decision 2026-11-06

- This repository uses ordinary Git and its own Make targets. CE, `ce task`,
  `ce access`, task_manager, and task-manager are not prerequisites and must not
  be invoked for development, worktree lifecycle, task records, verification,
  integration, or cleanup. This product-specific user decision overrides
  inherited instructions requiring those tools or proposing their adoption.
- Use an isolated task branch/worktree; preserve existing work and commits.
  Use Git directly for status, commit, task-branch push, integration and cleanup.
- Run the repository's applicable source checks and keep failures visible.
  Task-board state and card schemas are not source-quality gates. Documentation
  and workflow-only changes require diff checks and relevant Markdown/configuration
  or workflow checks; do not require unrelated Go test suites for prose edits.
- GitLab CI is not required to build, test, or contribute. GitHub Actions and
  server-side branch rules are separate and must still be respected where present.
- Existing `tasks/` documents and historical receipts are retained as records,
  not as an executable task-manager database or a reason to block source work.
- Preserve in-file locked regions. Do not change machine-global installed tools
  or unrelated repositories as part of this repository's workflow.

This file provides LLM-optimized guidance for Claude Code when working with this repository.

---

## Project Context

**Module**: `github.com/gizzahub/gzh-cli-core`
**Type**: Shared library for gzh-cli-* tools
**Consumer Go Version**: 1.26 (`go.mod` directive)
**Preferred Development / CI Toolchain**: 1.27.3 (main module directive / CI pin; not a consumer floor)

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
