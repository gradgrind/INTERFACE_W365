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
	w365 := ReadXML(f365)
	idmap := makeIdMap(&w365)

	/*
		coursemap := map[string]Course{}
		for _, c := range w365.Courses {
			coursemap[c.Id] = c
		}
		ecoursemap := map[string]EpochPlanCourse{}
		for _, c := range w365.EpochPlanCourses {
			ecoursemap[c.Id] = c
		}
		for i, d := range w365.Lessons {
			_, ok := coursemap[d.Course]
			if ok {
				if d.Fixed {
					fmt.Printf("*--- %02d: %+v\n", i, d)
				}
			} else {
				_, ok = ecoursemap[d.Course]
				if ok {
					fmt.Printf("*+++ %02d: %+v\n", i, d)
				} else {
					fmt.Printf("*::: %02d: %+v\n", i, d)
				}
			}
		}
	*/

	test_ids_exist(&w365, idmap)

	for _, n := range w365.Classes {
		read_divisions(idmap, n.Id)
	}
}
