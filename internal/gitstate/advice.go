// internal/advice/advice.go
package advice

import "github.com/danny-molnar/gitorch/internal/gitstate"

type Advice struct {
	Summary  string
	Messages []string
}

func Explain(s gitstate.State) Advice {
	var msgs []string

	if !s.HasRemote {
		msgs = append(msgs, "No remote configured (expected '"+s.RemoteName+"').")
		return Advice{
			Summary:  "Remote not found.",
			Messages: msgs,
		}
	}

	if s.BehindBy > 0 {
		msgs = append(msgs,
			"Your local "+s.DefaultBranch+" is "+itoa(s.BehindBy)+" commit(s) behind "+
				s.RemoteName+"/"+s.DefaultBranch+".",
			"Suggested: git pull --rebase "+s.RemoteName+" "+s.DefaultBranch,
		)
	}

	if s.AheadBy > 0 {
		msgs = append(msgs,
			"Your local "+s.DefaultBranch+" is "+itoa(s.AheadBy)+" commit(s) ahead of "+
				s.RemoteName+"/"+s.DefaultBranch+".",
			"Suggested: git push "+s.RemoteName+" "+s.DefaultBranch,
		)
	}

	if s.Dirty {
		msgs = append(msgs,
			"Working tree has uncommitted changes.",
			"Suggested: git status; commit or stash before updating.",
		)
	}

	if len(msgs) == 0 {
		msgs = append(msgs, "Local "+s.DefaultBranch+" is up to date with "+s.RemoteName+"/"+s.DefaultBranch+".")
	}

	return Advice{
		Summary:  msgs[0],
		Messages: msgs,
	}
}
