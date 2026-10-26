# Changelog

Release history for the gzh-cli-core Go library. Published versions are
identified by this repository's version tags and GitHub Releases.

## [0.1.0] - 2026-10-26

### Added

- Initial public library packages: `cli`, `config`, `errors`, `logger`,
  `testutil`, and `version`.
- Consumer minimum Go version 1.26.
- Release-note procedure and prospective API compatibility baseline.

### Compatibility

- The API compatibility promise starts only if `v0.1.0` is published from its
  tagged source. Development pseudo-versions are snapshots and are not a
  substitute for that baseline.
- There is no migration from an earlier tagged release.

### Known issues and qualification limits

- Existing pseudo-version consumers should verify their caller integration;
  there is no earlier stable tag from which to claim a migration guarantee.
- Compatibility execution used Go 1.26.8 on macOS arm64. Go 1.26.0 itself and
  other platform runtimes were not executed; cross-builds are compile checks.

### Verification

- Library implementation at `b9966055c92a232fe2044989e10b8fa825ae04a6` passed
  native `make check`, `make fmt-check build vet`, and race coverage:
  492/578 statements (85.121107%, unchanged 85% floor).
- [CI 40250886345](https://github.com/Gizzahub/gzh-cli-core/actions/runs/40250886345)
  passed the candidate's tests, lint and vulnerability checks.
- Go 1.26.8 fixed-source race tests and an external three-library consumer
  smoke test passed. CGO-disabled package cross-builds passed for macOS,
  Linux and Windows on amd64/arm64; no standalone executable is shipped.
- Publication requires the existing workflow to verify the actual tagged
  revision and extract these version-specific notes before creating a release.
