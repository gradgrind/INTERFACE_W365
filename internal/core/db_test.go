package core

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
	"path/filepath"
	"testing"
)

func TestFet(t *testing.T) {
	w365file := "../_testdata/fms_w365.json"
	// w365file := "../_testdata/test1.json"
	abspath, err := filepath.Abs(w365file)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", abspath)
	}
	data := w365tt.LoadJSON(abspath)

	db := DbTopLevel{}
	for _, d := range data.Days {
		db.Days = append(db.Days, Day{
			Id: Ref(d.Id), Name: d.Name, Tag: d.Tag})
	}
	fmt.Printf("  --> %+v\n", db)
}
