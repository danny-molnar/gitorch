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
}

// Inspect inspects the Git repository at repoPath and returns a summary of the
// state of DefaultBranch relative to RemoteName/DefaultBranch.
func Inspect(repoPath, defaultBranch, remoteName string) (State, error) {
	s := State{
		RepoPath:      repoPath,
		DefaultBranch: defaultBranch,
		RemoteName:    remoteName,
	}

	// Verify Git repository
	if err := runGit(repoPath, "rev-parse"); err != nil {
		return s, fmt.Errorf("not inside a Git repository: %w", err)
	}

	// Determine current branch (best effort; not strictly required for now)
	if out, err := runGitOutput(repoPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		s.CurrentBranch = strings.TrimSpace(out)
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

	// Determine ahead/behind status relative to remote/defaultBranch
	// git rev-list --left-right --count A...B
	// A = local, B = remote; output: "<ahead> <behind>"
	remoteRef := fmt.Sprintf("%s/%s", remoteName, defaultBranch)
	spec := fmt.Sprintf("%s...%s", defaultBranch, remoteRef)

	out, err := runGitOutput(repoPath, "rev-list", "--left-right", "--count", spec)
	if err != nil {
		return s, fmt.Errorf("unable to determine branch status: %w", err)
	}

	parts := strings.Fields(out)
	if len(parts) != 2 {
		return s, fmt.Errorf("unexpected rev-list output %q", out)
	}

	ahead, behind, err := parseAheadBehind(parts[0], parts[1])
	if err != nil {
		return s, err
	}
	s.AheadBy = ahead
	s.BehindBy = behind

	// Check if working tree is dirty
	if dirtyOut, err := runGitOutput(repoPath, "status", "--porcelain"); err == nil {
		s.Dirty = strings.TrimSpace(dirtyOut) != ""
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
