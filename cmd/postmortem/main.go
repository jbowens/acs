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

	counties, err := postmortem.ReadCounties(acsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading acs data: %s", err.Error())
		os.Exit(1)
	}

	// Just print the counties to stdout for now.
	for _, county := range counties {
		fmt.Printf("%s — %s — %s\n", county.ID, county.State, county.Name)
	}
}
