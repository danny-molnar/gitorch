// internal/advice/advice_test.go
package advice_test

import (
	"testing"

	"github.com/danny-molnar/gitorch/internal/advice"
	"github.com/danny-molnar/gitorch/internal/gitstate"
)

func hasCode(codes []advice.Code, want advice.Code) bool {
	for _, c := range codes {
		if c == want {
			return true
		}
	}
	return false
}

func TestExplain_UpToDate(t *testing.T) {
	s := gitstate.State{
		AheadBy:       0,
		BehindBy:      0,
		Dirty:         false,
		HasRemote:     true,
		DefaultBranch: "main",
		RemoteName:    "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeUpToDate) {
		t.Fatalf("expected CodeUpToDate, got %v", a.Codes)
	}
	if a.Summary == "" {
		t.Fatalf("expected non-empty Summary")
	}
}

func TestExplain_AheadOnly(t *testing.T) {
	s := gitstate.State{
		AheadBy:       3,
		BehindBy:      0,
		Dirty:         false,
		HasRemote:     true,
		DefaultBranch: "main",
		RemoteName:    "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeAheadOnly) {
		t.Fatalf("expected CodeAheadOnly")
	}
}

func TestExplain_BehindOnly(t *testing.T) {
	s := gitstate.State{
		AheadBy:       0,
		BehindBy:      2,
		HasRemote:     true,
		DefaultBranch: "main",
		RemoteName:    "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeBehindOnly) {
		t.Fatalf("expected CodeBehindOnly")
	}
}

func TestExplain_Diverged(t *testing.T) {
	s := gitstate.State{
		AheadBy:       2,
		BehindBy:      5,
		HasRemote:     true,
		DefaultBranch: "main",
		RemoteName:    "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeDiverged) {
		t.Fatalf("expected CodeDiverged")
	}
}

func TestExplain_Dirty(t *testing.T) {
	s := gitstate.State{
		Dirty: true,
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeDirty) {
		t.Fatalf("expected CodeDirty")
	}
}

func TestExplain_NoRemote(t *testing.T) {
	s := gitstate.State{
		HasRemote:     false,
		RemoteName:    "origin",
		DefaultBranch: "main",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeNoRemote) {
		t.Fatalf("expected CodeNoRemote")
	}
	if len(a.Details) == 0 {
		t.Fatalf("expected at least one detail")
	}
}

func TestExplain_DetachedHEAD(t *testing.T) {
	s := gitstate.State{
		DetachedHEAD:  true,
		DefaultBranch: "main",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeDetachedHEAD) {
		t.Fatalf("expected CodeDetachedHEAD")
	}
}

func TestExplain_NoUpstream(t *testing.T) {
	s := gitstate.State{
		HasRemote:     true,
		HasUpstream:   false,
		CurrentBranch: "feature/foo",
		RemoteName:    "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeNoUpstream) {
		t.Fatalf("expected CodeNoUpstream")
	}
}

func TestExplain_MissingDefaultLocal(t *testing.T) {
	s := gitstate.State{
		HasRemote:            true,
		MissingDefaultLocal:  true,
		MissingDefaultRemote: false,
		DefaultBranch:        "develop",
		RemoteName:           "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeMissingDefaultBranch) {
		t.Fatalf("expected CodeMissingDefaultBranch")
	}
	if a.Summary == "" {
		t.Fatalf("expected non-empty Summary for missing local default branch")
	}
}

func TestExplain_MissingDefaultRemote(t *testing.T) {
	s := gitstate.State{
		HasRemote:            true,
		MissingDefaultLocal:  false,
		MissingDefaultRemote: true,
		DefaultBranch:        "develop",
		RemoteName:           "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeMissingDefaultBranch) {
		t.Fatalf("expected CodeMissingDefaultBranch")
	}
}

func TestExplain_MissingDefaultBoth(t *testing.T) {
	s := gitstate.State{
		HasRemote:            true,
		MissingDefaultLocal:  true,
		MissingDefaultRemote: true,
		DefaultBranch:        "develop",
		RemoteName:           "origin",
	}

	a := advice.Explain(s)

	if !hasCode(a.Codes, advice.CodeMissingDefaultBranch) {
		t.Fatalf("expected CodeMissingDefaultBranch")
	}
}
