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
	geographyStateIdx     = 1
	geographyTypeIdx      = 2
	geographyComponentIdx = 3
	geographyRecNoIdx     = 4
	geographyIDIdx        = 48
	geographyNameIdx      = 49

	geographyTypeCounty    = "050"
	geographyTypeState     = "040"
	geographyComponentNone = "00"
)

// ImportStates reads all of the states out of the American
// Community Survey (ACS) geography files in the provided directory.
func ImportStates(acsPath string) ([]*State, error) {
	states := map[string]*State{}

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
			recNo, err := strconv.Atoi(rec[geographyRecNoIdx])
			if err != nil {
				return fmt.Errorf("invalid rec no: %q for geo %s", rec[geographyRecNoIdx], rec[geographyIDIdx])
			}
			if rec[geographyTypeIdx] == geographyTypeState && rec[geographyComponentIdx] == geographyComponentNone {
				state := rec[geographyStateIdx]
				states[state] = &State{
					ID:     rec[geographyIDIdx],
					Abbrev: rec[geographyStateIdx],
					Name:   rec[geographyNameIdx],
					RecNo:  recNo,
				}
			}
			if rec[geographyTypeIdx] == geographyTypeCounty {
				state := rec[geographyStateIdx]
				states[state].Counties = append(states[state].Counties, &County{
					ID:    rec[geographyIDIdx],
					State: rec[geographyStateIdx],
					Name:  strings.SplitN(rec[geographyNameIdx], ",", 2)[0],
					RecNo: recNo,
				})
			}
		}
		return nil
	})

	var statesSlice []*State
	for _, state := range states {
		statesSlice = append(statesSlice, state)
	}
	return statesSlice, err
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
		typ:    TotalPopulation{},
		tbl:    "B01003",
		seq:    "0003",
		offset: 130,
		count:  1,
	},
	{
		typ:    FoodStamps{},
		tbl:    "C22001",
		seq:    "0094",
		offset: 128,
		count:  3,
	},
}

// ImportACS imports all supported American Community Survey (ACS) statistics
// for the provided geographies.
func ImportACS(dir string, geographies []Geography) (map[string]*ACSStatistics, error) {
	allResults, err := importACS(dir, geographies, sequenceMappings...)
	if err != nil {
		return nil, err
	}

	m := make(map[string]*ACSStatistics, len(geographies))
	for id, results := range allResults {
		stats := new(ACSStatistics)
		for _, res := range results {
			switch v := res.(type) {
			case *FoodStamps:
				stats.FoodStamps = v
			case *TotalPopulation:
				stats.TotalPopulation = v
			default:
				return nil, fmt.Errorf("unexpected %T", res)
			}
		}
		m[id] = stats
	}
	return m, nil
}

// importACS imports the provided data tables for the provided geos.
// It returns a map from GeoID to a list of the hydrated table structs.
func importACS(dir string, geographies []Geography, dataTables ...dataTable) (map[string][]interface{}, error) {
	results := map[string][]interface{}{}
	byState := map[string][]Geography{}
	for _, g := range geographies {
		stateID := strings.ToLower(g.StateID())
		byState[stateID] = append(byState[stateID], g)
		results[g.GeoID()] = make([]interface{}, 0, len(dataTables))
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

			// For each geo, extract the geo's record and hydrate the relevant types.
			for _, geo := range stateCounties {
				rec := records[geo.RecordNo()]
				for _, tbl := range tables {
					v, err := parseSequence(tbl, rec[tbl.offset-1:(tbl.offset-1+tbl.count)])
					if err != nil {
						return nil, err
					}
					results[geo.GeoID()] = append(results[geo.GeoID()], v)
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
