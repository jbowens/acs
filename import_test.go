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
		{ID: "05000US44003", State: "RI", Name: "Kent County", RecNo: "0000008"},
		{ID: "05000US44005", State: "RI", Name: "Newport County", RecNo: "0000009"},
		{ID: "05000US44007", State: "RI", Name: "Providence County", RecNo: "0000010"},
		{ID: "05000US44009", State: "RI", Name: "Washington County", RecNo: "0000011"},
	}
	if !reflect.DeepEqual(counties, want) {
		t.Errorf("ReadCounties(): got:\n%#v\nwant:\n%#v", counties, want)
	}
}
