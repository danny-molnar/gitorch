// internal/gitstate/gitstate_test.go
package gitstate_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/danny-molnar/gitorch/internal/gitstate"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}

func newTempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	runGit(t, dir, "init", ".")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial")

	return dir
}

func cloneRepo(t *testing.T, remoteDir string) string {
	t.Helper()
	dir := t.TempDir()
	runGit(t, dir, "clone", remoteDir, ".")
	return dir
}

func bytesTrimSpace(b []byte) []byte {
	start, end := 0, len(b)
	for start < end && (b[start] == ' ' || b[start] == '\n' || b[start] == '\t' || b[start] == '\r') {
		start++
	}
	for end > start && (b[end-1] == ' ' || b[end-1] == '\n' || b[end-1] == '\t' || b[end-1] == '\r') {
		end--
	}
	return b[start:end]
}

func TestInspect_UpToDate(t *testing.T) {
	repo := newTempRepo(t)

	// assume default branch is whatever HEAD is currently on ("master" or "main")
	out, err := exec.Command("git", "-C", repo, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	defaultBranch := string(bytesTrimSpace(out))

	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare", ".")

	runGit(t, repo, "remote", "add", "origin", remoteDir)
	runGit(t, repo, "push", "-u", "origin", defaultBranch)

	state, err := gitstate.Inspect(repo, defaultBranch, "origin", "default")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.HasRemote {
		t.Fatalf("expected HasRemote=true")
	}
	if state.AheadBy != 0 || state.BehindBy != 0 {
		t.Fatalf("expected ahead/behind 0, got %d/%d", state.AheadBy, state.BehindBy)
	}
	if state.Dirty {
		t.Fatalf("expected clean working tree")
	}
	if state.MissingDefaultLocal || state.MissingDefaultRemote {
		t.Fatalf("did not expect missing default branch flags")
	}
}

func TestInspect_Dirty(t *testing.T) {
	repo := newTempRepo(t)

	out, err := exec.Command("git", "-C", repo, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	defaultBranch := string(bytesTrimSpace(out))

	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare", ".")
	runGit(t, repo, "remote", "add", "origin", remoteDir)
	runGit(t, repo, "push", "-u", "origin", defaultBranch)

	// Modify file without committing
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("changed"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	state, err := gitstate.Inspect(repo, defaultBranch, "origin", "default")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.Dirty {
		t.Fatalf("expected Dirty=true")
	}
}

func TestInspect_MissingDefaultLocal(t *testing.T) {
	repo := newTempRepo(t)

	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare", ".")
	runGit(t, repo, "remote", "add", "origin", remoteDir)

	// Push current branch so remote exists, but we will ask for a non-existent default branch
	out, err := exec.Command("git", "-C", repo, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	currentBranch := string(bytesTrimSpace(out))

	runGit(t, repo, "push", "-u", "origin", currentBranch)

	state, err := gitstate.Inspect(repo, "nonexistent", "origin", "default")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.MissingDefaultLocal {
		t.Fatalf("expected MissingDefaultLocal=true")
	}
	// remote default also missing, so this may be true as well; up to your implementation
}

func TestInspect_MissingDefaultRemote(t *testing.T) {
	repo := newTempRepo(t)

	// current branch is our default branch for the test
	out, err := exec.Command("git", "-C", repo, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	defaultBranch := string(bytesTrimSpace(out))

	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare", ".")
	runGit(t, repo, "remote", "add", "origin", remoteDir)

	// Deliberately do NOT push defaultBranch, so remote default branch doesn't exist
	state, err := gitstate.Inspect(repo, defaultBranch, "origin", "default")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.MissingDefaultRemote {
		t.Fatalf("expected MissingDefaultRemote=true")
	}
}

func TestInspect_DetachedHEAD(t *testing.T) {
	repo := newTempRepo(t)

	// Create a second commit so we can detach to the first one
	if err := os.WriteFile(filepath.Join(repo, "extra.txt"), []byte("extra"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runGit(t, repo, "add", "extra.txt")
	runGit(t, repo, "commit", "-m", "second")

	// Detach HEAD to the first commit
	out, err := exec.Command("git", "-C", repo, "rev-list", "--max-parents=0", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-list failed: %v\n%s", err, string(out))
	}
	firstCommit := string(bytesTrimSpace(out))
	runGit(t, repo, "checkout", firstCommit)

	state, err := gitstate.Inspect(repo, "main", "origin", "current") // branch/remote values don't matter much here
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.DetachedHEAD {
		t.Fatalf("expected DetachedHEAD=true")
	}
}

func TestInspect_CurrentMode_UpToDate(t *testing.T) {
	origin := t.TempDir()
	runGit(t, origin, "init", "--bare", ".")

	// Seed origin with an initial commit
	seed := newTempRepo(t)
	out, err := exec.Command("git", "-C", seed, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	defaultBranch := string(bytesTrimSpace(out))

	runGit(t, seed, "remote", "add", "origin", origin)
	runGit(t, seed, "push", "-u", "origin", defaultBranch)

	// Clone into working repo (will have upstream configured)
	repo := cloneRepo(t, origin)

	state, err := gitstate.Inspect(repo, defaultBranch, "origin", "current")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if state.DetachedHEAD {
		t.Fatalf("did not expect detached HEAD")
	}
	if !state.HasUpstream {
		t.Fatalf("expected HasUpstream=true")
	}
	if state.AheadBy != 0 || state.BehindBy != 0 {
		t.Fatalf("expected ahead/behind 0, got %d/%d", state.AheadBy, state.BehindBy)
	}
}

func TestInspect_CurrentMode_NoUpstream(t *testing.T) {
	repo := newTempRepo(t)

	// Create a new branch that has no upstream
	runGit(t, repo, "checkout", "-b", "feature/no-upstream")

	state, err := gitstate.Inspect(repo, "main", "origin", "current")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if state.DetachedHEAD {
		t.Fatalf("did not expect detached HEAD")
	}
	if state.HasUpstream {
		t.Fatalf("expected HasUpstream=false")
	}
	if state.AheadBy != 0 || state.BehindBy != 0 {
		t.Fatalf("expected ahead/behind 0 when no upstream, got %d/%d", state.AheadBy, state.BehindBy)
	}
}

func TestInspect_CurrentMode_BehindUpstream(t *testing.T) {
	origin := t.TempDir()
	runGit(t, origin, "init", "--bare", ".")

	// Seed origin
	seed := newTempRepo(t)
	out, err := exec.Command("git", "-C", seed, "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		t.Fatalf("git rev-parse failed: %v\n%s", err, string(out))
	}
	defaultBranch := string(bytesTrimSpace(out))

	runGit(t, seed, "remote", "add", "origin", origin)
	runGit(t, seed, "push", "-u", "origin", defaultBranch)

	// Two clones: one will be "behind"
	repoA := cloneRepo(t, origin) // will remain behind
	repoB := cloneRepo(t, origin) // will push an extra commit

	// Make a new commit in repoB and push
	if err := os.WriteFile(filepath.Join(repoB, "new.txt"), []byte("new"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runGit(t, repoB, "add", "new.txt")
	runGit(t, repoB, "commit", "-m", "new commit")
	runGit(t, repoB, "push", "origin", defaultBranch)

	// Now repoA should be behind after fetch (Inspect will fetch if HasRemote)
	state, err := gitstate.Inspect(repoA, defaultBranch, "origin", "current")
	if err != nil {
		t.Fatalf("Inspect failed: %v", err)
	}

	if !state.HasUpstream {
		t.Fatalf("expected HasUpstream=true")
	}
	if state.BehindBy == 0 {
		t.Fatalf("expected BehindBy>0, got %d (AheadBy=%d)", state.BehindBy, state.AheadBy)
	}
}
