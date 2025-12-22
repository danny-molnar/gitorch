# CHANGELOG.md

## v0.2.0

### Added

* New comparison mode for inspecting the current branch against its upstream
* `-mode` flag to switch between:
  * `current` (current branch vs upstream, default)
  * `default` (default branch vs remote default)
* Improved handling for detached HEAD and missing upstream configuration

### Changed

* Default behaviour now compares the current branch rather than the default branch
* Ahead/behind calculation refactored into a shared helper
* Dirty working tree detection now runs independently of comparison mode


## v0.1.0

Initial release.

### Features

* Inspect default branch vs remote default branch
* Detect ahead / behind / diverged states
* Detect dirty working tree
* Detect missing remote
* Detect missing local or remote default branch
* Detect detached HEAD
* Detect missing upstream for current branch
* Provide human-readable guidance
* Support JSON output for tooling integration
* Support quiet (summary-only) output

### Internal

* Clear separation between git state inspection and advice generation
* Comprehensive unit tests for gitstate and advice layers
