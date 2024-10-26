package fet

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
)

// "Atomic Groups" are needed especially for the class handling.
// They should only be built for divisions which have lessons.
// So first the Lessons must be consulted for their Courses
// and thus their groups – which can then be marked. Finally the divisions
// can be filtered on the basis of these marked groups.

func makeAtomicGroups(fetinfo *fetInfo) {
	// Mark the Groups used by Lessons.
	markedGroups := map[db.DbRef]bool{}
	for _, l := range fetinfo.db.Lessons {
		lc := l.Course
		cix, ok := fetinfo.courses[lc]
		if ok {
			// It is a normal course.
			for _, g := range fetinfo.db.Courses[cix].Groups {
				markedGroups[g] = true
			}
		} else {
			_, ok = fetinfo.supercourses[lc]
			if !ok {
				msg := fmt.Sprintf("#BUG# Lesson %d has invalid course.", l.Id)
				panic(msg)
			}
			// It is a supercourse, go throught its subcourses.
			for _, sub := range fetinfo.supersubs[lc] {
				subix, ok := fetinfo.subcourses[sub]
				if !ok {
					msg := fmt.Sprintf("#BUG# subcourses[%d].", sub)
					panic(msg)
				}
				for _, g := range fetinfo.db.SubCourses[subix].Groups {
					markedGroups[g] = true
				}
			}
		}
	}
	// Go through the classes inspecting their Divisions.
	for _, cl := range fetinfo.db.Classes {
		ags := []string{}

		divs := [][]db.DbRef{}
		for _, d := range cl.Divisions {
			dok := false
			for _, g := range d.Groups {
				if markedGroups[g] {
					dok = true
					break
				}
			}
			if dok {
				if len(ags) == 0 {
					ags = []string{cl.Tag}
				}
				divs = append(divs, d.Groups)
				agsx := []string{}
				for _, ag := range ags {
					for _, g := range d.Groups {
						agsx = append(agsx, ag+
							"~"+fetinfo.ref2grouponly[g])
					}
				}
				ags = agsx
			}
		}
		fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		fmt.Printf("     --> %+v\n", ags)
	}
}
