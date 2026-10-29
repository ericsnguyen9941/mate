// revieweval renders the production review brief for a frozen evaluation case.
// It never spawns a harness; run.py owns the explicitly requested live pilot.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nguyenngocanh94/mate/internal/mateassets"
	"github.com/nguyenngocanh94/mate/internal/store"
)

func main() {
	report := flag.String("check-report", "", "check a passing report against the production contract")
	prior := flag.Int("prior-count", -1, "also check this number of prior finding closures")
	flag.Parse()
	if *report != "" {
		data, err := os.ReadFile(*report)
		if err == nil {
			err = store.ValidateReviewPass(string(data))
		}
		if err == nil && *prior >= 0 {
			err = store.ValidatePriorFindingClosures(string(data), *prior)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("pass contract satisfied")
		return
	}
	var p mateassets.ReviewParams
	if err := json.NewDecoder(os.Stdin).Decode(&p); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(mateassets.RenderReviewBrief(p))
}
