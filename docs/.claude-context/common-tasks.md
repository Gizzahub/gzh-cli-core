# Common Tasks - gzh-cli-core

## Adding a New Utility Package

### Pattern

```text
gzh-cli-core/<pkg>/
├── doc.go          # package doc comment
├── <name>.go       # implementation
├── <name>_test.go  # tests (80%+ coverage)
└── errors.go       # sentinel errors (if needed)
```

### Steps

1. Create the package directory under `gzh-cli-core/`
2. Add `doc.go` with package-level GoDoc comment
3. Implement functions — keep them stateless where possible
4. Write table-driven tests using `testing` + `testutil`
5. Add sentinel errors to `errors/errors.go` if the package defines new error conditions

### Example: env helper

```go
// gzh-cli-core/config/env.go
func GetEnv(key string, prefix ...string) string {
    p := DefaultEnvPrefix // "GZH"
    if len(prefix) > 0 {
        p = prefix[0]
    }
    if p != "" {
        key = p + "_" + key
    }
    return os.Getenv(key)
}
```

________________________________________________________________________

## Using the Config Package

### Environment variables

All env lookups go through `config.GetEnv` which prepends `GZH_`:

```go
val := config.GetEnv("CONFIG_DIR")   // reads GZH_CONFIG_DIR
```

Override prefix for project-specific tools:

```go
val := config.GetEnv("API_KEY", "MYTOOL")  // reads MYTOOL_API_KEY
```

### Typed helpers

| Function                     | Returns            | Fallback           |
| ---------------------------- | ------------------ | ------------------ |
| `GetEnvOr(key, def, ...)`    | `string`           | `def` if empty     |
| `GetEnvBool(key, ...)`       | `bool`             | `false` if unset   |
| `GetEnvBoolOr(key, def, ...)`| `bool`             | `def` if unset     |
| `GetEnvInt(key, ...)`        | `(int, bool)`      | `(0, false)`       |
| `GetEnvIntOr(key, def, ...)` | `int`              | `def` if unset     |
| `GetEnvDuration(key, ...)`   | `(time.Duration, bool)` | `(0, false)` |
| `GetEnvList(key, ...)`       | `[]string`         | `nil` if unset     |

### Config directory

```go
dir := config.GetConfigDirectory()    // GZH_CONFIG_DIR or ~/.config/gzh-manager
err := config.EnsureConfigDirectory() // mkdir -p with 0o750
```

### YAML/TOML config loading

```go
loader := config.NewLoader("mytool")
cfg, err := loader.Load("config.yaml")  // searches CWD, config dir, home
```

________________________________________________________________________

## Using the Logger

### Basic usage

```go
import "github.com/gizzahub/gzh-cli-core/logger"

log := logger.New(os.Stderr, logger.LevelInfo)
log.Info("starting server", "port", 8080)
log.Error("connection failed", "err", err)
```

### Levels

| Constant       | Integer | String   |
| -------------- | ------- | -------- |
| `LevelDebug`   | 0       | `DEBUG`  |
| `LevelInfo`    | 1       | `INFO`   |
| `LevelWarn`    | 2       | `WARN`   |
| `LevelError`   | 3       | `ERROR`  |

### Global logger

```go
logger.SetGlobal(log)
logger.Global().Info("shared message")
```

________________________________________________________________________

## Using the Errors Package

### Sentinel errors

```go
import "github.com/gizzahub/gzh-cli-core/errors"

if errors.Is(err, errors.ErrNotFound) {
    // handle not-found
}
```

| Sentinel            | Meaning                  |
| ------------------- | ------------------------ |
| `ErrNotFound`       | Resource not found       |
| `ErrInvalidInput`   | Bad user input           |
| `ErrConfigNotFound` | Config file missing      |
| `ErrInvalidConfig`  | Config parse/validation  |
| `ErrUnauthorized`   | Auth failed              |
| `ErrTimeout`        | Operation timed out      |
| `ErrPermission`     | Permission denied        |
| `ErrAlreadyExists`  | Duplicate create         |
| `ErrNotSupported`   | Platform/feature unsup.  |

### Wrapping

```go
err := errors.Wrap(err, errors.ErrNotFound)
err := errors.Wrapf(err, errors.ErrInvalidConfig, "missing key %s", key)
```

### Validation helpers

```go
errors.ValidateRequired("name", name)          // nil if non-empty
errors.ValidateInRange("port", port, 1, 65535) // nil if in range
```

________________________________________________________________________

## Using the CLI Helpers

### Root command

```go
import "github.com/gizzahub/gzh-cli-core/cli"

cmd := cli.NewRootCmd(cli.RootConfig{
    Name:    "mytool",
    Short:   "Does something useful",
    Version: "1.0.0",
})
```

### Output formatting

```go
out := cli.NewOutput()
out.SetFormat("json")  // text, json, yaml

out.Print(map[string]any{"status": "ok"})
```

### LLM-friendly formatter

```go
out.SetFormat("llm")  // compact, token-efficient output
```

________________________________________________________________________

## Using TestUtil

### Temp files

```go
func TestSomething(t *testing.T) {
    dir := testutil.TempDir(t)                    // auto-cleanup
    file := testutil.TempFile(t, "config.yaml", "key: value")
}
```

### Assertions

```go
testutil.AssertEqual(t, expected, actual)
testutil.AssertNoError(t, err)
testutil.AssertErrorContains(t, err, "not found")
```

### Environment capture

```go
func TestEnvSensitive(t *testing.T) {
    restore := testutil.CaptureEnv(t)
    defer restore()
    os.Setenv("MY_VAR", "test")
}
```

________________________________________________________________________

## Using Version

```go
import "github.com/gizzahub/gzh-cli-core/version"

info := version.Get()
fmt.Printf("%s %s (built %s)\n", info.Version, info.GitCommit, info.BuildDate)
```

Set at build time:

```bash
go build -ldflags "-X github.com/gizzahub/gzh-cli-core/version.Version=1.0.0"
```
