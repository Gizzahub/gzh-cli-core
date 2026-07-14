# Architecture Guide - gzh-cli-core

## Package Map

```text
gzh-cli-core/
├── config/     # Environment + config file loading
│   ├── env.go       # GetEnv, GetEnvOr, GetEnvBool, GetEnvInt, ...
│   ├── dir.go       # GetConfigDirectory, EnsureConfigDirectory
│   └── loader.go    # YAML/TOML config file loader
├── cli/        # Cobra command helpers
│   ├── root.go      # NewRootCmd(RootConfig)
│   ├── output.go    # Output (text/json/yaml/llm)
│   ├── flags.go     # GlobalFlags
│   └── llm_formatter.go  # LLM-friendly output
├── errors/     # Sentinel errors + validation
│   ├── errors.go         # 9 sentinels + Wrap/Wrapf
│   └── validation.go     # ValidateRequired, ValidateInRange
├── logger/     # Structured logging
│   ├── logger.go    # Logger, SimpleLogger, Level
│   └── global.go    # SetGlobal/Global
├── testutil/   # Test helpers
│   ├── testutil.go  # TempDir, TempFile, Capture
│   ├── assert.go    # AssertEqual, AssertNoError
│   └── env.go       # CaptureEnv, Chdir
└── version/    # Build-time version info
    └── version.go   # Info{Version,GitCommit,BuildDate,GoVersion,Platform}
```

## Dependency Direction

gzh-cli-core has **zero internal dependencies** — each package is standalone.

```
errors    → (stdlib only)
config    → (stdlib only)
version   → (stdlib only)
logger    → (stdlib only)
cli       → cobra, yaml.v3
testutil  → testing (test-only)
```

All gzh-cli-* sibling projects depend on core:

```
gzh-cli          → gzh-cli-core
gzh-cli-gitforge → gzh-cli-core
gzh-cli-net-env  → gzh-cli-core  (config.GetConfigDirectory)
gzh-cli-os-env   → gzh-cli-core
...
```

## Design Principles

1. **Stateless functions preferred** — `GetEnv`, `GetConfigDirectory` are pure reads
2. **Sentinel errors** — stable identity for `errors.Is()` across packages
3. **No business logic** — core provides utilities, not features
4. **Minimal dependencies** — only cobra + yaml.v3 for cli package
5. **Test-friendly** — `testutil` package exists for all sibling projects

## Adding New Shared Utilities

When multiple sibling projects need the same helper:

1. Check if it belongs in core (generic, no project-specific logic)
2. Create the package under `gzh-cli-core/`
3. Add tests (80%+ coverage)
4. Export clean API — no internal types leaked
5. Update consuming projects' `go.mod`

## Consumer Pattern

```go
import (
    "github.com/gizzahub/gzh-cli-core/config"
    "github.com/gizzahub/gzh-cli-core/errors"
    "github.com/gizzahub/gzh-cli-core/logger"
)

func main() {
    log := logger.New(os.Stderr, logger.LevelInfo)
    dir := config.GetConfigDirectory()

    if err := setup(dir); err != nil {
        log.Error("setup failed", "err", err)
        if errors.Is(err, errors.ErrConfigNotFound) {
            log.Info("run 'mytool init' to create config")
        }
        os.Exit(1)
    }
}
```
