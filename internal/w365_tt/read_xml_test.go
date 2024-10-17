package w365_tt

import (
	"fmt"
	"log"
	"strings"
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

	idmap := map[string]string{}
	for _, c := range w365.Days {
		idmap[c.Id] = "Day"
	}
	for _, c := range w365.Hours {
		idmap[c.Id] = "Hour"
	}
	for _, c := range w365.Absences {
		idmap[c.Id] = "Absence"
	}
	for _, c := range w365.Teachers {
		idmap[c.Id] = "Teacher"
	}
	for _, c := range w365.Subjects {
		idmap[c.Id] = "Subject"
	}
	for _, c := range w365.Rooms {
		idmap[c.Id] = "Room"
	}
	for _, c := range w365.Classes {
		idmap[c.Id] = "Class"
	}
	for _, c := range w365.Groups {
		idmap[c.Id] = "Group"
	}
	for _, c := range w365.Divisions {
		idmap[c.Id] = "Division"
	}
	for _, c := range w365.Courses {
		idmap[c.Id] = "Course"
	}
	for _, c := range w365.EpochPlanCourses {
		idmap[c.Id] = "EpochPlanCourse"
	}
	for _, c := range w365.Lessons {
		idmap[c.Id] = "Lesson"
	}
	for _, c := range w365.Fractions {
		idmap[c.Id] = "Fraction"
	}

	for _, c := range w365.Teachers {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Teachers.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Teachers.Categories %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Rooms {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Rooms.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Rooms.Categories %s\n", x)
				}
			}
		}
		if len(c.RoomGroups) != 0 {
			for _, x := range strings.Split(c.RoomGroups, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Rooms.RoomGroups %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Classes {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Classes.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Classes.Categories %s\n", x)
				}
			}
		}
		if len(c.Divisions) != 0 {
			for _, x := range strings.Split(c.Divisions, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Classes.Divisions %s\n", x)
				}
			}
		}
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Classes.Groups %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Groups {
		if len(c.Absences) != 0 {
			for _, x := range strings.Split(c.Absences, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Groups.Absences %s\n", x)
				}
			}
		}
		if len(c.Categories) != 0 {
			for _, x := range strings.Split(c.Categories, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Groups.Categories %s\n", x)
				}
			}
		}
	}
	for _, c := range w365.Divisions {
		if len(c.Groups) != 0 {
			for _, x := range strings.Split(c.Groups, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Teachers.Groups %s\n", x)
				}
			}
		}
	}

	for _, c := range w365.Fractions {
		if len(c.SuperGroups) != 0 {
			for _, x := range strings.Split(c.SuperGroups, ",") {
				_, ok := idmap[x]
				if !ok {
					fmt.Printf("  !!! Fractions.SuperGroups %s\n", x)
				}
			}
		}
	}
}
