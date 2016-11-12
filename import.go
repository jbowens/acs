package postmortem

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

const (
	geographyStateIdx = 1
	geographyTypeIdx  = 2
	geographyRecNoIdx = 4
	geographyIDIdx    = 48
	geographyNameIdx  = 49

	geographyTypeCounty = "050"
)

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
		defer f.Close()
		recs, err := csv.NewReader(f).ReadAll()
		if err != nil {
			return err
		}
		for _, rec := range recs {
			if rec[geographyTypeIdx] != geographyTypeCounty {
				continue
			}
			recNo, err := strconv.Atoi(rec[geographyRecNoIdx])
			if err != nil {
				return fmt.Errorf("invalid rec no: %q for geo %s", rec[geographyRecNoIdx], rec[geographyIDIdx])
			}
			counties = append(counties, County{
				ID:    rec[geographyIDIdx],
				State: rec[geographyStateIdx],
				Name:  strings.SplitN(rec[geographyNameIdx], ",", 2)[0],
				RecNo: recNo,
			})
		}
		return nil
	})
	return counties, err
}

type dataTable struct {
	typ    interface{}
	tbl    string
	seq    string
	offset int
	count  int
}

var sequenceMappings = []dataTable{
	{
		typ:    FoodStamps{},
		tbl:    "C22001",
		seq:    "0094",
		offset: 128,
		count:  3,
	},
}

// importACS imports the provided data tables for the provided counties.
// It returns a map from GeoID to a list of the hydrated table structs.
func importACS(dir string, counties []County, dataTables ...dataTable) (map[string][]interface{}, error) {
	results := map[string][]interface{}{}
	byState := map[string][]County{}
	for _, c := range counties {
		stateID := strings.ToLower(c.State)
		byState[stateID] = append(byState[stateID], c)
		results[c.ID] = make([]interface{}, 0, len(dataTables))
	}

	tblsBySeq := map[string][]dataTable{}
	for _, tbl := range dataTables {
		tblsBySeq[tbl.seq] = append(tblsBySeq[tbl.seq], tbl)
	}

	for stateID, stateCounties := range byState {
		for seq, tables := range tblsBySeq {
			filename := filepath.Join(dir, fmt.Sprintf("e20151%s%s000.txt", stateID, seq))
			f, err := os.Open(filename)
			if err != nil {
				return nil, fmt.Errorf("error opening '%s': %s", filename, err)
			}

			// Read the entire file into memory.
			records, err := csv.NewReader(f).ReadAll()
			if err != nil {
				f.Close()
				return nil, err
			}
			err = f.Close()
			if err != nil {
				return nil, err
			}

			// For each county, extract the county's record and hydrate the relevant types.
			for _, county := range stateCounties {
				rec := records[county.RecNo]
				for _, tbl := range tables {
					v, err := parseSequence(tbl, rec[tbl.offset-1:(tbl.offset-1+tbl.count)])
					if err != nil {
						return nil, err
					}
					results[county.ID] = append(results[county.ID], v)
				}
			}
		}
	}
	return results, nil
}

func parseSequence(tbl dataTable, rec []string) (interface{}, error) {
	if tbl.count != len(rec) {
		fmt.Errorf("table expects %d fields, got %d", tbl.count, len(rec))
	}
	typ := reflect.TypeOf(tbl.typ)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type; got %T", tbl.typ)
	}
	if typ.NumField() != tbl.count {
		return nil, fmt.Errorf("struct has %d fields, table expects %d", typ.NumField(), tbl.count)
	}

	v := reflect.New(typ)
	structV := v.Elem()
	for i := 0; i < tbl.count; i++ {
		num, err := strconv.Atoi(rec[i])
		if err != nil {
			return nil, err
		}
		structV.Field(i).SetInt(int64(num))
	}
	return v.Interface(), nil
}
