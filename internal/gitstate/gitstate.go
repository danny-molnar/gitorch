// internal/gitstate/gitstate.go
package gitstate

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type State struct {
	RepoPath      string
	CurrentBranch string

	DefaultBranch string
	RemoteName    string

	AheadBy   int
	BehindBy  int
	Dirty     bool
	HasRemote bool

	DetachedHEAD         bool
	HasUpstream          bool
	Upstream             string
	MissingDefaultLocal  bool
	MissingDefaultRemote bool
}

func aheadBehind(repoPath, leftRef, rightRef string) (ahead int, behind int, err error) {
	spec := fmt.Sprintf("%s...%s", leftRef, rightRef)

	out, err := runGitOutput(repoPath, "rev-list", "--left-right", "--count", spec)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to determine branch status: %w", err)
	}

	parts := strings.Fields(out)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected rev-list output %q", out)
	}

	return parseAheadBehind(parts[0], parts[1])
}

// Inspect inspects the Git repository at repoPath and returns a summary of the
// state of DefaultBranch relative to RemoteName/DefaultBranch.
func Inspect(repoPath, defaultBranch, remoteName string, mode string) (State, error) {
	s := State{
		RepoPath:      repoPath,
		DefaultBranch: defaultBranch,
		RemoteName:    remoteName,
	}

	// Verify Git repository
	if err := runGit(repoPath, "rev-parse"); err != nil {
		return s, fmt.Errorf("not inside a Git repository: %w", err)
	}

	// Determine current branch
	if out, err := runGitOutput(repoPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		br := strings.TrimSpace(out)
		if br == "HEAD" {
			s.DetachedHEAD = true
		} else {
			s.CurrentBranch = br
		}
	}

	// Determine upstream for current branch (if not detached)
	if !s.DetachedHEAD && s.CurrentBranch != "" {
		if up, err := runGitOutput(repoPath,
			"rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}",
		); err == nil {
			s.HasUpstream = true
			s.Upstream = strings.TrimSpace(up)
		} else {
			s.HasUpstream = false
		}
	}

	// Check for remote
	if out, err := runGitOutput(repoPath, "remote"); err == nil {
		for _, r := range strings.Split(out, "\n") {
			if strings.TrimSpace(r) == remoteName {
				s.HasRemote = true
				break
			}
		}
	}

	// Fetch latest remote information (only if remote exists)
	if s.HasRemote {
		if err := runGit(repoPath, "fetch", "--quiet", remoteName); err != nil {
			return s, fmt.Errorf("unable to fetch latest information from remote %q: %w", remoteName, err)
		}
	}

	// Check if working tree is dirty
	if dirtyOut, err := runGitOutput(repoPath, "status", "--porcelain"); err == nil {
		s.Dirty = strings.TrimSpace(dirtyOut) != ""
	}

	// Determine ahead/behind based on mode.
	switch mode {
	case "default":
		// Determine ahead/behind status relative to remote/defaultBranch.
		// First detect whether the local and remote default branches actually exist.

		localRef := fmt.Sprintf("refs/heads/%s", defaultBranch)
		remoteRef := fmt.Sprintf("refs/remotes/%s/%s", remoteName, defaultBranch)

		if err := runGit(repoPath, "show-ref", "--verify", "--quiet", localRef); err != nil {
			s.MissingDefaultLocal = true
			return s, nil
		}

		// In default mode, remote must exist; otherwise we cannot compare to remote/defaultBranch.
		if !s.HasRemote {
			s.MissingDefaultRemote = true
			return s, nil
		}

		if err := runGit(repoPath, "show-ref", "--verify", "--quiet", remoteRef); err != nil {
			s.MissingDefaultRemote = true
			return s, nil
		}

		// At this point both local and remote default branches exist.
		ahead, behind, err := aheadBehind(
			repoPath,
			defaultBranch,
			fmt.Sprintf("%s/%s", remoteName, defaultBranch),
		)
		if err != nil {
			return s, err
		}
		s.AheadBy = ahead
		s.BehindBy = behind

	case "current":
		// Compare current branch against its upstream.
		if s.DetachedHEAD {
			return s, nil
		}
		if s.CurrentBranch == "" {
			return s, nil
		}
		if !s.HasUpstream || strings.TrimSpace(s.Upstream) == "" {
			return s, nil
		}

		// Use HEAD vs upstream for comparison.
		ahead, behind, err := aheadBehind(repoPath, "HEAD", s.Upstream)
		if err != nil {
			return s, err
		}
		s.AheadBy = ahead
		s.BehindBy = behind

	default:
		return s, fmt.Errorf("invalid mode %q (expected: current or default)", mode)
	}

	return s, nil
}

func runGit(repoPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	if repoPath != "" && repoPath != "." {
		cmd.Dir = repoPath
	}
	return cmd.Run()
}

func runGitOutput(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if repoPath != "" && repoPath != "." {
		cmd.Dir = repoPath
	}
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v failed: %w (output: %s)", args, err, buf.String())
	}
	return buf.String(), nil
}

func parseAheadBehind(aheadStr, behindStr string) (int, int, error) {
	ahead, err := parseInt(aheadStr)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse ahead count %q: %w", aheadStr, err)
	}
	behind, err := parseInt(behindStr)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse behind count %q: %w", behindStr, err)
	}
	return ahead, behind, nil
}

func parseInt(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
