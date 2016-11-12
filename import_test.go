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
		{ID: "05000US44003", State: "RI", Name: "Kent County"},
		{ID: "05000US44005", State: "RI", Name: "Newport County"},
		{ID: "05000US44007", State: "RI", Name: "Providence County"},
		{ID: "05000US44009", State: "RI", Name: "Washington County"},
	}
	if !reflect.DeepEqual(counties, want) {
		t.Errorf("ReadCounties(): got:\n%#v\nwant:\n%#v", counties, want)
	}
}
