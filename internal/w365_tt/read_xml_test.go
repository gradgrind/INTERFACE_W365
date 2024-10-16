package w365_tt

import (
	"fmt"
	"log"
	"testing"

	"github.com/ncruces/zenity"
)

func TestReadXML(t *testing.T) {
	fmt.Println("\n############## TestReadXML")
	const defaultPath = "../_testdata/*.xml"
	f365, err := zenity.SelectFile(
		zenity.Filename(defaultPath),
		zenity.FileFilter{
			Name:     "Waldorf-365 TT-export",
			Patterns: []string{"*.xml"},
			CaseFold: false,
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n ***** Reading %s *****\n", f365)

	ReadXML(f365)

}
