package postmortem

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
)

const (
	geographyStateIdx = 1
	geographyTypeIdx  = 2
	geographyIDIdx    = 48
	geographyNameIdx  = 49

	geographyTypeCounty = "050"
)

// County represents an individual United States county.
type County struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
}

// ReadCounties reads all of the counties out of the American
// Community Survey (ACS) geography files in the provided directory.
func ReadCounties(acsPath string) ([]County, error) {
	var counties []County

	err := filepath.Walk(acsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".csv" {
			return nil // skip non csvs
		}
		if info.Name()[0] != 'g' {
			return nil // skip non-geo files
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		recs, err := csv.NewReader(f).ReadAll()
		if err != nil {
			return err
		}
		for _, rec := range recs {
			if rec[geographyTypeIdx] != geographyTypeCounty {
				continue
			}
			counties = append(counties, County{
				ID:    rec[geographyIDIdx],
				State: rec[geographyStateIdx],
				Name:  strings.SplitN(rec[geographyNameIdx], ",", 2)[0],
			})
		}
		return nil
	})
	return counties, err
}
