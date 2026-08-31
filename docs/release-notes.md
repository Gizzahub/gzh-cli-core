# Release-note procedure

This procedure prepares release notes. It does not authorize a tag, a hosted
release, or a release workflow change.

## Maintain the Unreleased section

Keep [`CHANGELOG.md`](../CHANGELOG.md) under an `## [Unreleased]` heading until
publication is separately approved. Record only user-visible changes in these
categories when they apply:

- Added
- Changed
- Deprecated
- Removed
- Fixed
- Security

The section must state changes to the public API, consumer minimum Go version,
configuration or wire contracts, known issues, and verification evidence.

## Before publication

Prepare a release note from the following template. Leave the version and date
unresolved until tag and publication approval is given.

```markdown
## [<version>] - <YYYY-MM-DD>

### Summary

<user-visible summary>

### Added / Changed / Fixed / Security

<grouped entries from Unreleased>

### Compatibility and migration

<state whether public API, Go minimum, config, or wire contracts changed>
<for every breaking change, show before/after use and a migration path>

### Known issues

<known limitations and their workarounds>

### Verification

<exact revision and completed local/hosted verification evidence>
```

Breaking changes require an explicit before/after migration section. Removing a
deprecated API is still a breaking change. Do not claim a release has been
published until the separately authorized tag and release steps are complete.
