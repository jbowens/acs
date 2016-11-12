package postmortem

import (
	"reflect"
	"testing"
)

func TestReadCounties(t *testing.T) {
	counties, err := ReadCounties("test_data/acsgeos")
	if err != nil {
		t.Fatal(err)
	}
	want := []County{
		{ID: "05000US44003", State: "RI", Name: "Kent County", RecNo: 8},
		{ID: "05000US44005", State: "RI", Name: "Newport County", RecNo: 9},
		{ID: "05000US44007", State: "RI", Name: "Providence County", RecNo: 10},
		{ID: "05000US44009", State: "RI", Name: "Washington County", RecNo: 11},
	}
	if !reflect.DeepEqual(counties, want) {
		t.Errorf("ReadCounties(): got:\n%#v\nwant:\n%#v", counties, want)
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
