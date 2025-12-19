// cmd/gitorch/main.go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/danny-molnar/gitorch/internal/advice"
	"github.com/danny-molnar/gitorch/internal/gitstate"
)

func main() {
	defaultBranch := flag.String("branch", "main", "default branch to compare against")
	remoteName := flag.String("remote", "origin", "remote name to compare against")

	quiet := flag.Bool("quiet", false, "print summary only")
	jsonOut := flag.Bool("json", false, "print JSON output")

	flag.Parse()

	state, err := gitstate.Inspect(".", *defaultBranch, *remoteName)
	if err != nil {
		log.Fatal(err)
	}

	adv := advice.Explain(state)

	// JSON mode wins over quiet/human modes
	if *jsonOut {
		if err := printJSON(os.Stdout, state, adv); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Quiet: just the one-line summary
	if *quiet {
		fmt.Println(adv.Summary)
		return
	}

	// Default: human-readable summary + details
	fmt.Println(adv.Summary)

	if len(adv.Details) > 0 {
		fmt.Println()
		for _, msg := range adv.Details {
			fmt.Println("•", msg)
		}
	}
}

type jsonOutput struct {
	State  gitstate.State `json:"state"`
	Advice advice.Advice  `json:"advice"`
}

func printJSON(w *os.File, s gitstate.State, a advice.Advice) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(jsonOutput{
		State:  s,
		Advice: a,
	})
}
