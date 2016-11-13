package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jbowens/postmortem"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "need path to ACS data")
		os.Exit(1)
	}
	acsPath := args[0]

	counties, err := postmortem.ImportCounties(acsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading acs geo data: %s", err.Error())
		os.Exit(1)
	}

	results, err := postmortem.ImportACS(acsPath, counties)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading acs estimate data: %s", err.Error())
		os.Exit(1)
	}

	// Just print the state population estimates to stdout for now.
	statePopulations := map[string]int{}
	for county, stats := range results {
		statePopulations[county.State] = statePopulations[county.State] + stats.TotalPopulation.Total
	}
	for state, pop := range statePopulations {
		fmt.Printf("%s — %d\n", state, pop)
	}
}
