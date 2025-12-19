# README.md

# gitorch

**gitorch** is a small CLI that inspects your Git repository and tells you what
your branch needs next — pull, push, rebase, or nothing at all.

It focuses on *state and guidance*, not raw Git output.

---

## What gitorch does

Given a repository, a default branch, and a remote, gitorch:

* Compares the local default branch against the remote default branch
* Detects:

  * ahead / behind / diverged states
  * dirty working tree
  * missing remote
  * missing default branch (local and/or remote)
  * detached HEAD
  * missing upstream for the current branch
* Produces:

  * a one-line summary
  * optional detailed guidance
  * optional machine-readable JSON

---

## Installation

```bash
go install github.com/danny-molnar/gitorch/cmd/gitorch@latest
```

---

## Usage

```bash
gitorch [flags]
```

### Flags

| Flag      | Description                        | Default  |
| --------- | ---------------------------------- | -------- |
| `-branch` | Default branch to compare against  | `main`   |
| `-remote` | Remote name to compare against     | `origin` |
| `-quiet`  | Print summary only                 | `false`  |
| `-json`   | Print JSON output (state + advice) | `false`  |

> Note: gitorch uses Go’s standard `flag` package, so flags are single-dash
> (e.g. `-quiet`, not `--quiet`).

---

## Examples

### Up-to-date branch

```bash
$ gitorch
Branch is up to date with remote.
```

---

### Branch ahead of remote

```bash
$ gitorch
Branch is ahead of remote.

• Branch is ahead of origin/main by 2 commits; consider pushing.
```

---

### Branch behind remote

```bash
$ gitorch
Branch is behind remote.

• Branch is behind origin/main by 3 commits; consider pulling.
```

---

### Diverged branch

```bash
$ gitorch
Branch has diverged from remote.

• Branch has diverged from origin/main (ahead by 2, behind by 5); consider `git pull --rebase`.
```

---

### Dirty working tree

```bash
$ gitorch
Branch is up to date with remote.

• Working tree has uncommitted changes.
```

---

### No remote configured

```bash
$ gitorch
No remote "origin" configured; add it with `git remote add origin <url>`.
```

---

### Detached HEAD

```bash
$ gitorch
HEAD is detached.

• HEAD is detached; checkout a branch (e.g. `git switch main`) before running gitorch.
```

---

### Quiet mode

```bash
$ gitorch -quiet
Branch is behind remote.
```

---

### JSON output

```bash
$ gitorch -json
```

Example output:

```json
{
  "state": {
    "repoPath": ".",
    "currentBranch": "main",
    "defaultBranch": "main",
    "remoteName": "origin",
    "aheadBy": 0,
    "behindBy": 2,
    "dirty": false,
    "hasRemote": true
  },
  "advice": {
    "codes": ["behind_only"],
    "summary": "Branch is behind remote.",
    "details": [
      "Branch is behind origin/main by 2 commits; consider pulling."
    ]
  }
}
```

---

## Design goals

* Be explicit, not clever
* Prefer guidance over raw Git output
* Fail gracefully when information is missing
* Stay dependency-light and scriptable

---

## License

MIT
