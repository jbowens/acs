package postmortem

import (
	"reflect"
	"testing"
)

func TestImportStates(t *testing.T) {
	states, err := ImportStates("test_data/acsgeos")
	if err != nil {
		t.Fatal(err)
	}
	want := []*State{
		&State{
			ID:     "04000US44",
			Abbrev: "RI",
			Name:   "Rhode Island",
			RecNo:  1,
			Counties: []*County{
				{ID: "05000US44003", State: "RI", Name: "Kent County", RecNo: 8},
				{ID: "05000US44005", State: "RI", Name: "Newport County", RecNo: 9},
				{ID: "05000US44007", State: "RI", Name: "Providence County", RecNo: 10},
				{ID: "05000US44009", State: "RI", Name: "Washington County", RecNo: 11},
			},
		},
	}
	if !reflect.DeepEqual(states, want) {
		t.Errorf("ImportStates(): got:\n%#v\nwant:\n%#v", states, want)
	}
}

func TestParseSequence(t *testing.T) {
	exampleTbl := dataTable{
		typ:    FoodStamps{},
		tbl:    "C22001",
		seq:    "0094",
		offset: 0,
		count:  3,
	}

	v, err := parseSequence(exampleTbl, []string{"5", "3", "2"})
	if err != nil {
		t.Fatal(err)
	}
	want := &FoodStamps{Total: 5, Yes: 3, No: 2}
	if !reflect.DeepEqual(v, want) {
		t.Errorf("got=%#v want=%#v", v, want)
	}
}

func TestImportACS(t *testing.T) {
	foodStampsTbl := dataTable{
		typ:    FoodStamps{},
		tbl:    "C22001",
		seq:    "0094",
		offset: 128,
		count:  3,
	}

	prov := &County{ID: "05000US44007", State: "RI", Name: "Providence County", RecNo: 10}
	results, err := importACS("test_data/acsri", []Geography{prov}, foodStampsTbl)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]interface{}{
		"05000US44007": []interface{}{&FoodStamps{Total: 49555, Yes: 3977, No: 45578}},
	}
	if !reflect.DeepEqual(results, want) {
		t.Errorf("got=%#v, want=%#v", results, want)
	}
}
