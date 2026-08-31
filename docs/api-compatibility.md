# API compatibility policy

## Scope and start of the baseline

No `v0.1.0` tag exists yet. Pseudo-versions before that tag are development
snapshots and are not covered by this compatibility promise. If `v0.1.0` is
published, every exported declaration in that tagged source becomes the initial
public API baseline for `github.com/gizzahub/gzh-cli-core`.

The module's `go 1.26` directive is the consumer minimum Go version. The
`toolchain go1.26.7` directive and the CI pin are preferred development and CI
settings, not a consumer compatibility floor.

## Public packages

| Package | Compatibility surface |
| --- | --- |
| `cli` | Exported flag/configuration types, `Output`, root/version helpers, package output helpers, and their documented selector/flag behavior. `Output.Print` routes `json`, `yaml`/`yml`, and `llm`; other selector values use the shipped text fallback. |
| `config` | `DefaultEnvPrefix`, environment and config-directory helpers, `Loader` and its methods, `DefaultPaths`, and `Save`. |
| `errors` | The nine exported sentinel error identities plus exported wrapping, inspection, and validation helpers. |
| `logger` | `Level` constants, `Logger`, `SimpleLogger`, `NopLogger`, constructors, methods, and global helpers. |
| `testutil` | Exported temp, capture, environment, working-directory, and assertion helpers whose behavior is documented or tested. |
| `version` | Build variables, `Info`, its exported fields and JSON/YAML tags, its methods, `Get`, and `LdFlags`. |

The tagged source, GoDoc, and these package-level boundaries are authoritative;
this document intentionally does not duplicate every identifier.

## Compatibility rules

Although this is a `v0` module, this project keeps the stricter product policy:
breaking changes to the published `v0.1.0` baseline require a SemVer major
release and a migration note.

The following require compatibility review and, when breaking, a migration note:

- removing, renaming, or changing an exported name, type, signature, method
  set, constant, variable, or exported struct field;
- adding a method to `logger.Logger`, because external consumers can implement
  that interface;
- adding an exported struct field when it can break consumers that use unkeyed
  literals;
- changing `version.Info` JSON/YAML tags, sentinel error identities, or the
  documented wrapping behavior of exported error helpers;
- changing documented CLI flags/defaults, environment prefix/fallback behavior,
  or documented helper behavior; and
- removing a deprecated API.

Before the first tag there is no earlier released API to migrate from. Once a
tagged baseline exists, every breaking release note must show the prior use,
the replacement, and the migration steps.

## Exclusions

The following are not compatibility commitments unless a later public document
states otherwise:

- pre-tag pseudo-versions and unexported implementation details, including
  `cli.llmFormatter` and exact internal formatting structure;
- exact dependency versions, repository layout, Make targets, CI mechanics,
  and the preferred development/CI toolchain;
- undocumented logger timestamp/context iteration details; and
- `testutil.Capture*` panic or large-output behavior and `AssertNil` or
  `AssertNotNil` behavior for non-nil-able values.

If a product decision is needed to make `llm` experimental or to give `table`
its own non-fallback formatter, make that decision and add tests before a tag.
This policy otherwise preserves the selector behavior that ships in the tagged
source without promising the unexported formatter implementation.
