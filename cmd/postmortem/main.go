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

	// Just print the counties to stdout for now.
	for county, stats := range results {
		foodStampsPct := 100.0 * (float64(stats.FoodStamps.Yes) / float64(stats.FoodStamps.Total))
		fmt.Printf("%s (%s) — %.2f\n", county.Name, county.State, foodStampsPct)
	}
}
