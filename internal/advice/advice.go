// internal/advice/advice.go
package advice

import (
	"fmt"

	"github.com/danny-molnar/gitorch/internal/gitstate"
)

type Advice struct {
	Summary  string
	Messages []string
}

func Explain(s gitstate.State) Advice {
	var msgs []string

	if s.DefaultBranch == "" {
		s.DefaultBranch = "main"
	}
	if s.RemoteName == "" {
		s.RemoteName = "origin"
	}

	if !s.HasRemote {
		msgs = append(msgs,
			fmt.Sprintf("No remote named %q found for this repository.", s.RemoteName),
			"Add a remote or run gitorch with the correct -remote flag.",
		)
		return Advice{
			Summary:  msgs[0],
			Messages: msgs,
		}
	}

	remoteRef := fmt.Sprintf("%s/%s", s.RemoteName, s.DefaultBranch)

	// Summary & messages based on behind/ahead/dirty
	switch {
	case s.BehindBy > 0 && s.AheadBy == 0:
		msgs = append(msgs,
			fmt.Sprintf("Your local %q branch is behind %q by %d commit(s).", s.DefaultBranch, remoteRef, s.BehindBy),
			fmt.Sprintf("Suggested: git pull --rebase %s %s", s.RemoteName, s.DefaultBranch),
			"Alternatively: git merge if you prefer a merge-based workflow.",
		)
	case s.BehindBy == 0 && s.AheadBy > 0:
		msgs = append(msgs,
			fmt.Sprintf("Your local %q branch is ahead of %q by %d commit(s).", s.DefaultBranch, remoteRef, s.AheadBy),
			fmt.Sprintf("Suggested: git push %s %s", s.RemoteName, s.DefaultBranch),
		)
	case s.BehindBy > 0 && s.AheadBy > 0:
		msgs = append(msgs,
			fmt.Sprintf("Your local %q branch has diverged from %q (ahead by %d, behind by %d).", s.DefaultBranch, remoteRef, s.AheadBy, s.BehindBy),
			fmt.Sprintf("Suggested: git pull --rebase %s %s", s.RemoteName, s.DefaultBranch),
			"Review the history before pushing to avoid surprises.",
		)
	default:
		msgs = append(msgs,
			fmt.Sprintf("Your local %q branch is up-to-date with %q.", s.DefaultBranch, remoteRef),
		)
	}

	if s.Dirty {
		msgs = append(msgs,
			"Working tree has uncommitted or unstaged changes.",
			"Suggested: git status; commit or stash before rebasing, merging, or pushing.",
		)
	}

	return Advice{
		Summary:  msgs[0],
		Messages: msgs,
	}
}
