package main

import (
	"flag"
	"fmt"
	"math"
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

	sum := 0
	stateReps := map[string]int{}
	statePopulations := map[string]int{}
	for county, stats := range results {
		if county.State == "PR" || county.State == "DC" {
			continue
		}
		sum += stats.TotalPopulation.Total
		statePopulations[county.State] = statePopulations[county.State] + stats.TotalPopulation.Total
	}

	representatives := 435

	// Every state gets at least 1 representative.
	for state := range statePopulations {
		stateReps[state]++
		representatives--
	}
	for ; representatives > 0; representatives-- {
		var highestPriority float64
		var highestPriorityState string
		for state, n := range stateReps {
			priority := float64(statePopulations[state]) / math.Sqrt(float64(n)*(float64(n)+1))
			if priority > highestPriority {
				highestPriority, highestPriorityState = priority, state
			}
		}
		stateReps[highestPriorityState]++
	}

	// Print the resulting repsentative counts.
	for state, reps := range stateReps {
		fmt.Printf("%02d\t%s\n", reps, state)
	}
}
