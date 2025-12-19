// internal/advice/advice.go
package advice

import (
	"fmt"

	"github.com/danny-molnar/gitorch/internal/gitstate"
)

type Code string

const (
	CodeUpToDate   Code = "up_to_date"
	CodeAheadOnly  Code = "ahead_only"
	CodeBehindOnly Code = "behind_only"
	CodeDiverged   Code = "diverged"
	CodeDirty      Code = "dirty"
	CodeNoRemote   Code = "no_remote"

	CodeDetachedHEAD         Code = "detached_head"
	CodeNoUpstream           Code = "no_upstream"
	CodeMissingDefaultBranch Code = "missing_default_branch"
)

type Advice struct {
	Codes   []Code
	Summary string
	Details []string
}

func Explain(s gitstate.State) Advice {
	var a Advice

	if !s.HasRemote {
		a.Codes = append(a.Codes, CodeNoRemote)

		remote := s.RemoteName
		if remote == "" {
			remote = "origin"
		}

		msg := fmt.Sprintf(
			"No remote %q configured; add it with `git remote add %s <url>`.",
			remote, remote,
		)
		a.Details = append(a.Details, msg)

		if a.Summary == "" {
			a.Summary = msg
		}
	}

	// Detached HEAD
	if s.DetachedHEAD {
		a.Codes = append(a.Codes, CodeDetachedHEAD)
		msg := "HEAD is detached; checkout a branch (e.g. `git switch " + s.DefaultBranch + "`) before running gitorch."
		a.Details = append(a.Details, msg)
		if a.Summary == "" {
			a.Summary = "HEAD is detached."
		}
	}

	// No upstream for current branch
	if !s.DetachedHEAD && s.CurrentBranch != "" && !s.HasUpstream && s.HasRemote {
		a.Codes = append(a.Codes, CodeNoUpstream)

		remote := s.RemoteName
		if remote == "" {
			remote = "origin"
		}

		msg := fmt.Sprintf(
			"Branch %q has no upstream configured; set it with `git push -u %s %s`.",
			s.CurrentBranch, remote, s.CurrentBranch,
		)
		a.Details = append(a.Details, msg)
		if a.Summary == "" {
			a.Summary = "Current branch has no upstream configured."
		}
	}

	// Missing default branch (local, remote, or both)
	if s.MissingDefaultLocal || s.MissingDefaultRemote {
		a.Codes = append(a.Codes, CodeMissingDefaultBranch)

		switch {
		case s.MissingDefaultLocal && s.MissingDefaultRemote:
			msg := fmt.Sprintf("Default branch %q does not exist locally or on %s.", s.DefaultBranch, s.RemoteName)
			a.Details = append(a.Details, msg)
			if a.Summary == "" {
				a.Summary = msg
			}
		case s.MissingDefaultLocal:
			msg := fmt.Sprintf("Local default branch %q does not exist; fetch or create it before comparing.", s.DefaultBranch)
			a.Details = append(a.Details, msg)
			if a.Summary == "" {
				a.Summary = msg
			}
		case s.MissingDefaultRemote:
			msg := fmt.Sprintf("Remote default branch %s/%s does not exist; push or configure the remote default branch.", s.RemoteName, s.DefaultBranch)
			a.Details = append(a.Details, msg)
			if a.Summary == "" {
				a.Summary = msg
			}
		}
	}

	if s.Dirty {
		a.Codes = append(a.Codes, CodeDirty)
		a.Details = append(a.Details, "Working tree has uncommitted changes.")
	}

	switch {
	case s.AheadBy == 0 && s.BehindBy == 0:
		a.Codes = append(a.Codes, CodeUpToDate)
		if a.Summary == "" {
			a.Summary = "Branch is up to date with remote."
		}

	case s.AheadBy > 0 && s.BehindBy == 0:
		a.Codes = append(a.Codes, CodeAheadOnly)
		a.Details = append(a.Details,
			fmt.Sprintf("Branch is ahead of %s/%s by %d commits; consider pushing.",
				s.RemoteName, s.DefaultBranch, s.AheadBy),
		)
		if a.Summary == "" {
			a.Summary = "Branch is ahead of remote."
		}

	case s.AheadBy == 0 && s.BehindBy > 0:
		a.Codes = append(a.Codes, CodeBehindOnly)
		a.Details = append(a.Details,
			fmt.Sprintf("Branch is behind %s/%s by %d commits; consider pulling.",
				s.RemoteName, s.DefaultBranch, s.BehindBy),
		)
		if a.Summary == "" {
			a.Summary = "Branch is behind remote."
		}

	case s.AheadBy > 0 && s.BehindBy > 0:
		a.Codes = append(a.Codes, CodeDiverged)
		a.Details = append(a.Details,
			fmt.Sprintf("Branch has diverged from %s/%s (ahead by %d, behind by %d); consider `git pull --rebase`.",
				s.RemoteName, s.DefaultBranch, s.AheadBy, s.BehindBy),
		)
		if a.Summary == "" {
			a.Summary = "Branch has diverged from remote."
		}
	}

	if a.Summary == "" {
		a.Summary = "No specific advice."
	}

	return a
}
