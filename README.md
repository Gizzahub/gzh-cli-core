# gzh-cli-core

Shared core library for gzh-cli-* tools.

## Installation

```bash
go get github.com/gizzahub/gzh-cli-core
```

## Packages

| Package | Purpose |
|---------|---------|
| `logger` | Structured logging with levels and context |
| `testutil` | Testing utilities (temp files, assertions, env helpers) |
| `errors` | Error types and wrapping utilities |
| `config` | YAML configuration loading with env override |
| `cli` | Cobra command helpers and output formatting |
| `version` | Build version information |

## Usage

### Logger

```go
import "github.com/gizzahub/gzh-cli-core/logger"

// Create named logger
log := logger.New("myapp")
log.SetLevel(logger.LevelDebug)

log.Info("starting", "port", 8080)
log.Error("failed", "error", err)

// With context
reqLog := log.WithContext("request_id", "abc123")
reqLog.Info("processing request")

// Global logger
logger.SetDefault(log)
logger.Info("using global logger")
```

### TestUtil

```go
import "github.com/gizzahub/gzh-cli-core/testutil"

func TestExample(t *testing.T) {
    // Temp files
    dir := testutil.TempDir(t)
    file := testutil.TempFile(t, "test.txt", "content")

    // Assertions
    testutil.AssertNoError(t, err)
    testutil.AssertEqual(t, got, want)
    testutil.AssertContains(t, output, "expected")

    // Environment
    testutil.SetEnv(t, "MY_VAR", "value")
    testutil.Chdir(t, dir)

    // Output capture
    output := testutil.Capture(func() {
        fmt.Println("captured")
    })
}
```

### Errors

```go
import "github.com/gizzahub/gzh-cli-core/errors"

// Sentinel errors
if errors.Is(err, errors.ErrNotFound) {
    // handle not found
}

// Wrapping
err = errors.WrapOp("open file", err)
err = errors.WrapWithMessage(err, "additional context")

// Validation errors
return errors.RequiredFlag("output")
return errors.MutuallyExclusive("verbose", "quiet")
return errors.Range("port", 1, 65535)
```

### Config

```go
import "github.com/gizzahub/gzh-cli-core/config"

type AppConfig struct {
    Name string `yaml:"name"`
    Port int    `yaml:"port"`
}

// Load from default paths
loader := config.NewLoader("myapp")
var cfg AppConfig
if err := loader.LoadOrDefault(&cfg); err != nil {
    // handle error
}

// Environment variables
port := config.GetEnvIntOr("PORT", 8080)
debug := config.GetEnvBool("DEBUG")
timeout := config.GetEnvDurationOr("TIMEOUT", 30*time.Second)
```

### CLI

```go
import "github.com/gizzahub/gzh-cli-core/cli"

func main() {
    root := cli.NewRootCmd(cli.RootConfig{
        Name:    "myapp",
        Short:   "My application",
        Version: version.Get().Short(),
    })

    var flags cli.GlobalFlags
    cli.AddGlobalFlags(root, &flags)

    cli.Execute(root)
}

// Output helpers
cli.Success("Operation completed")
cli.Error("Operation failed: %v", err)
cli.Warning("Deprecated feature")
```

### Version

```go
import "github.com/gizzahub/gzh-cli-core/version"

// Get version info
info := version.Get()
fmt.Println(info.String())  // Full version info
fmt.Println(info.Short())   // Just version number
fmt.Println(info.Full())    // Version with git commit

// Build with ldflags
// go build -ldflags "-X github.com/gizzahub/gzh-cli-core/version.Version=1.0.0 ..."
```

## Development

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover
```

## License

MIT
