package w365_tt

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/base"
	"log"
	"testing"

	"github.com/ncruces/zenity"
)

func readfile() (W365TT, IdMap) {
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
	return w365, idmap
}

func TestReadXML(t *testing.T) {
	//w365, idmap := readfile()

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
}

func TestIdsExist(t *testing.T) {
	w365, idmap := readfile()
	test_ids_exist(&w365, idmap)
}

func TestClassDivisions(t *testing.T) {
	w365, idmap := readfile()
	for _, n := range w365.Classes {
		divs := read_divisions(idmap, n.Id)
		fmt.Printf("  ++ %d%s: %+v\n", n.Level, n.Letter, divs)
	}
}

func TestDivisionGroups(t *testing.T) {
	w365, idmap := readfile()
	used_groups := test_courses(&w365, idmap)
	for _, n := range w365.Classes {
		divs := read_divisions(idmap, n.Id)
		fmt.Printf("++++++ %s +++++++++++++++++++++++++++++++++\n", n.Tag())
		fmt.Printf("   Class Group used: %d\n", used_groups[n.Id])
		for _, cdiv := range divs {
			for _, g := range cdiv.Groups {
				if used_groups[g] == 0 {
					fmt.Printf(" --- %s: %s\n", cdiv.Name, g)
				}
			}
		}
		fmt.Println("---------------------------------------------")
	}
}

func TestFractions(t *testing.T) {
	w365, idmap := readfile()
	test_fractions(&w365, idmap)
}

func TestGroups(t *testing.T) {
	w365, idmap := readfile()
	test_lesson_groups(&w365, idmap)
}

func Test2DB(t *testing.T) {
	w365, idmap := readfile()
	db := collectData(&w365, idmap)
	base.SaveJSON(db.Records, "")
}
