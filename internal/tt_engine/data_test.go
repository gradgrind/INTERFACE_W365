package tt_engine

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/core"
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

func TestData(t *testing.T) {
	w365file := "../_testdata/fms_w365.json"
	// w365file := "../_testdata/test1.json"
	abspath, err := filepath.Abs(w365file)
	if err != nil {
		log.Fatalf("Couldn't resolve file path: %s\n", abspath)
	}

	stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	logpath := stempath + ".log"
	core.OpenLog(logpath)

	data := w365tt.LoadJSON(abspath)
	db := core.MoveDb(data)
	ttdata := initData(db)

	fmt.Printf("*** ResourceMap: %+v\n", ttdata.ResourceMap)
}
