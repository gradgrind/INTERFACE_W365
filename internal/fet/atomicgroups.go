package fet

import (
	"fmt"
	"strings"
)

// "Atomic Groups" are needed especially for the class handling.
// They should only be built for divisions which have lessons.
// So first the Lessons must be consulted for their Courses
// and thus their groups – which can then be marked. Finally the divisions
// can be filtered on the basis of these marked groups.

type AtomicGroup struct {
	Class  Ref
	Groups []Ref
	Tag    string
}

func makeAtomicGroups(fetinfo *fetInfo) {
	// Mark the Groups used by Lessons.
	markedGroups := map[Ref]bool{}
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
				msg := fmt.Sprintf("#BUG# Lesson %s has invalid course.", l.Id)
				panic(msg)
			}
			// It is a supercourse, go through its subcourses.
			for _, sub := range fetinfo.supersubs[lc] {
				subix, ok := fetinfo.subcourses[sub]
				if !ok {
					msg := fmt.Sprintf("#BUG# subcourses[%s].", sub)
					panic(msg)
				}
				for _, g := range fetinfo.db.SubCourses[subix].Groups {
					markedGroups[g] = true
				}
			}
		}
	}

	// An atomic group is an ordered list of single groups from each division.
	fetinfo.atomicgroups = map[Ref][]AtomicGroup{}
	// Go through the classes inspecting their Divisions. Retain only those
	// which have lessons.
	fetinfo.classdivisions = map[Ref][][]Ref{}
	for _, cl := range fetinfo.db.Classes {
		agi := [][]Ref{{}}
		divs := [][]Ref{}
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
				agix := [][]Ref{}

				for _, ag := range agi {
					for _, g := range d.Groups {
						gx := make([]Ref, len(ag)+1)
						copy(gx, append(ag, g))
						agix = append(agix, gx)
					}
				}
				agi = agix
			}
		}
		fetinfo.classdivisions[cl.Id] = divs
		//fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		//fmt.Printf("     --> %+v\n", agi)

		// Make AtomicGroups
		aglist := []AtomicGroup{}
		for _, ag := range agi {
			glist := []string{}
			for _, g := range ag {
				glist = append(glist, fetinfo.ref2grouponly[g])
			}
			ago := AtomicGroup{
				Class:  cl.Id,
				Groups: ag,
				Tag:    fmt.Sprintf("%s#%s", cl.Tag, strings.Join(glist, "~")),
			}
			aglist = append(aglist, ago)
		}
		//fmt.Printf("     ++> %+v\n", aglist)

		g2ags := map[Ref][]AtomicGroup{}
		//		xg2ags := map[string][]string{}
		i := len(divs)
		n := 1
		for i > 0 {
			i--
			a := 0 // ag index

			for a < len(aglist) {
				for _, g := range divs[i] {
					for j := 0; j < n; j++ {
						g2ags[g] = append(g2ags[g], aglist[a])
						//						xg2ags[fetinfo.ref2fet[g]] = append(
						//							xg2ags[fetinfo.ref2fet[g]], aglist[a])
						a++
					}
				}
			}

			n *= len(divs[i])
		}
		//fmt.Printf("     ++> %+v\n", xg2ags)
		if len(divs) != 0 {
			fetinfo.atomicgroups[cl.Id] = aglist
			for g, agl := range g2ags {
				agls := []string{}
				for _, ag := range agl {
					agls = append(agls, ag.Tag)
				}
				//fmt.Printf("     ++ %s: %+v\n", fetinfo.ref2fet[g], agls)
				fetinfo.atomicgroups[g] = agl
			}
		} else {
			fetinfo.atomicgroups[cl.Id] = []AtomicGroup{}
		}
	}
	//fmt.Println("\n +++++++++++++++++++++++++++")
	//printAtomicGroups(fetinfo)
}

func printAtomicGroups(fetinfo *fetInfo) {
	for _, cl := range fetinfo.db.Classes {
		agls := []string{}
		for _, ag := range fetinfo.atomicgroups[cl.Id] {
			agls = append(agls, ag.Tag)
		}
		fmt.Printf("  ++ %s: %+v\n", fetinfo.ref2fet[cl.Id], agls)
		for _, div := range fetinfo.classdivisions[cl.Id] {
			for _, g := range div {
				agls := []string{}
				for _, ag := range fetinfo.atomicgroups[g] {
					agls = append(agls, ag.Tag)
				}
				fmt.Printf("    -- %s: %+v\n", fetinfo.ref2fet[g], agls)
			}
		}
	}
}
