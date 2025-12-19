# CHANGELOG.md

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
