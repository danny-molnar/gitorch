// cmd/gitorch/main.go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/danny-molnar/gitorch/internal/advice"
	"github.com/danny-molnar/gitorch/internal/gitstate"
)

func main() {
	defaultBranch := flag.String("branch", "main", "default branch to compare against")
	remoteName := flag.String("remote", "origin", "remote name to compare against")
	flag.Parse()

	state, err := gitstate.Inspect(".", *defaultBranch, *remoteName)
	if err != nil {
		log.Fatal(err)
	}

	adv := advice.Explain(state)

	fmt.Println(adv.Summary)
	if len(adv.Messages) > 1 {
		fmt.Println()
		for _, msg := range adv.Messages[1:] {
			fmt.Println("•", msg)
		}
	}
}
