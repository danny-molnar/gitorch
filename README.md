# gitorch (formerly gitorchk)

**gitorch** is a lightweight command-line tool written in Go that helps you understand the state of your local Git repository at a glance.  
It inspects your branch, compares it with its remote, and provides clear, actionable advice — whether you're ahead, behind, diverged, or have uncommitted changes.

This tool grew from a simple behind/ahead checker into a structured CLI with internal Git state analysis and a human-friendly advice engine.

---

## Installation

Until the module/repo rename is complete, install using:

```bash
go install github.com/danny-molnar/gitorchk/cmd/gitorch@latest
```

Ensure your Go environment is configured correctly and your `$GOPATH/bin` or `$GGOBIN` is in your system PATH.

---

## Usage

Run `gitorch` inside any Git repository:

```bash
gitorch
```

By default, gitorch compares:

- local `main`  
- against `origin/main`

### Example output

```
Your local "main" branch is behind "origin/main" by 3 commit(s).

• Suggested: git pull --rebase origin main
• Alternatively: git merge if you prefer a merge-based workflow.
```

If your working tree is dirty:

```
Working tree has uncommitted or unstaged changes.
• Suggested: git status; commit or stash before rebasing, merging, or pushing.
```

---

## Flags

Override the default branch or remote:

```bash
gitorch -branch develop
gitorch -remote upstream
gitorch -branch release -remote github
```

Flags:

| Flag        | Description                                   | Default |
|-------------|-----------------------------------------------|---------|
| `-branch`   | Branch to compare against its upstream remote | `main`  |
| `-remote`   | Name of the Git remote                        | `origin` |

---

## Features

Current capabilities:

- Detects whether your local branch is **ahead**, **behind**, or **diverged** from its remote.
- Provides **actionable advice** (rebase, pull, push, etc.).
- Warns when the working tree is **dirty** (staged/unstaged changes).
- Supports custom branch and remote names.
- Structured architecture (`cmd/`, `internal/gitstate`, `internal/advice`) for future expansion.

Planned features (Issue #5):

- JSON output mode (`--json`)
- More detailed divergence analysis
- Improved CLI UX (quiet mode, summary-only mode)
- Module/repo rename to `gitorch`
- Test coverage for gitstate and advice modules

---

## Project Structure

```
.
├── cmd/
│   └── gitorch/
│       └── main.go         # CLI entrypoint
├── internal/
│   ├── gitstate/           # Git inspection logic
│       └── gitstate.go
│   └── advice/             # Advice engine for user-facing messaging
│       └── advice.go
├── go.mod
└── README.md
```

---

## Contributing

Contributions are welcome!  
If you have an idea, find a bug, or want to help shape gitorch’s direction:

- Open an issue  
- Submit a PR  
- Or join the discussion on Issue #5 (refactor & rebrand)

---

## License

This project is licensed under the MIT License — see `LICENSE.md` for details.
