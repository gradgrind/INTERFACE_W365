package fet

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
	"strings"
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
		ags := [][]string{{}}

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
				divs = append(divs, d.Groups)
				agsx := [][]string{}

				for _, ag := range ags {
					for _, g := range d.Groups {
						gx := append(ag, fetinfo.ref2grouponly[g])
						agsx = append(agsx, gx)
					}
				}
				ags = agsx
			}
		}
		fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		aglist := []string{}
		for _, ag := range ags {
			aglist = append(aglist, cl.Tag+"#"+strings.Join(ag, "/"))
		}
		fmt.Printf("     --> %+v\n", aglist)

		g2ags := map[db.DbRef][]int{}
		xg2ags := map[string][]string{}
		i := len(divs)
		n := 1
		for i > 0 {
			i--
			a := 0 // ag index

			for a < len(ags) {
				for _, g := range divs[i] {
					for j := 0; j < n; j++ {
						g2ags[g] = append(g2ags[g], a)
						xg2ags[fetinfo.ref2fet[g]] = append(
							xg2ags[fetinfo.ref2fet[g]], aglist[a])
						a++
					}
				}
			}

			n *= len(divs[i])
		}
		fmt.Printf("     ++> %+v\n", xg2ags)
	}
}
