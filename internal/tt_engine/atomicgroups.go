package tt_engine

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/core"
	"log"
	"strings"
)

// "Atomic Groups" are needed especially for the class handling.
// They should only be built for divisions which have lessons.
// So first the Lessons must be consulted for their Courses
// and thus their groups – which can then be marked. Finally the divisions
// can be filtered on the basis of these marked groups.

const CLASS_GROUP_SEP = "."
const ATOMIC_GROUP_SEP1 = "#"
const ATOMIC_GROUP_SEP2 = "~"

type AtomicGroup struct {
	Class  core.Ref
	Groups []core.Ref
	Tag    string
}

func filterDivisions(db *core.DbTopLevel) map[core.Ref][][]core.Ref {
	// Prepare filtered versions of the class Divisions containing only
	// those Divisions which have Groups used in Lessons.

	// Collect groups used in Lessons. Get them from the courses in the
	// db.CourseLessons map, which only includes courses with lessons.
	usedgroups := map[core.Ref]bool{}
	for cref := range db.CourseLessons {
		csc := db.Elements[cref]
		c, ok := csc.(*core.Course)
		if ok {
			for _, g := range c.Groups {
				usedgroups[g] = true
			}
		} else {
			for _, sc := range db.SuperSubs[cref] {
				for _, g := range sc.Groups {
					usedgroups[g] = true
				}
			}
		}
	}
	// Filter the class divisions, discarding the division names.
	cdivs := map[core.Ref][][]core.Ref{}
	for _, c := range db.Classes {
		divs := [][]core.Ref{}
		for _, div := range c.Divisions {
			for _, gref := range div.Groups {
				if usedgroups[gref] {
					divs = append(divs, div.Groups)
					break
				}
			}
		}
		cdivs[c.Id] = divs
	}
	return cdivs
}

func makeAtomicGroups(
	db *core.DbTopLevel,
	classDivs map[core.Ref][][]core.Ref,
) map[core.Ref][]AtomicGroup {
	// An atomic group is an ordered list of single groups from each division.
	atomicGroups := map[core.Ref][]AtomicGroup{}
	// Go through the classes inspecting their Divisions. Only those which
	// have Lessons are considered.
	// Build a list-basis for the atomic groups based on the Cartesian product.
	for _, cl := range db.Classes {
		divs, ok := classDivs[cl.Id]
		if !ok {
			log.Fatalf("*BUG* fetinfo.classDivisions[%s]\n", cl.Id)
		}
		// The atomic groups will be built as a list of lists of Refs.
		agrefs := [][]core.Ref{{}}
		for _, dglist := range divs {
			// Add another division – increases underlying list lengths.
			agrefsx := [][]core.Ref{}
			for _, ag := range agrefs {
				// Extend each of the old list items by appending each
				// group of the new division in turn – multiplies the
				// total number of atomic groups.
				for _, g := range dglist {
					gx := make([]core.Ref, len(ag)+1)
					copy(gx, append(ag, g))
					agrefsx = append(agrefsx, gx)
				}
			}
			agrefs = agrefsx
		}
		//fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		//fmt.Printf("     --> %+v\n", agrefs)

		// Make AtomicGroups
		aglist := []AtomicGroup{}
		for _, ag := range agrefs {
			glist := []string{}
			for _, gref := range ag {
				glist = append(glist, db.Elements[gref].(*core.Group).Tag)
			}
			ago := AtomicGroup{
				Class:  cl.Id,
				Groups: ag,
				Tag: cl.Tag + ATOMIC_GROUP_SEP1 +
					strings.Join(glist, ATOMIC_GROUP_SEP2),
			}
			aglist = append(aglist, ago)
		}

		// Map the individual groups to their atomic groups.
		g2ags := map[core.Ref][]AtomicGroup{}
		count := 1
		divIndex := len(divs)
		for divIndex > 0 {
			divIndex--
			divGroups := divs[divIndex]
			agi := 0 // ag index
			for agi < len(aglist) {
				for _, g := range divGroups {
					for j := 0; j < count; j++ {
						g2ags[g] = append(g2ags[g], aglist[agi])
						agi++
					}
				}
			}
			count *= len(divGroups)
		}
		if len(divs) != 0 {
			atomicGroups[cl.Id] = aglist
			for g, agl := range g2ags {
				atomicGroups[g] = agl
			}
		} else {
			atomicGroups[cl.Id] = []AtomicGroup{}
		}
	}
	return atomicGroups
}

func classOrGroup(db *core.DbTopLevel, ref core.Ref) string {
	// Return the tag for the Class or Group, in the latter case including
	// the Class.
	ginfo, ok := db.GroupInfoMap[ref]
	if ok {
		ctag := db.Elements[ginfo.Class].(*core.Class).Tag
		return ctag + CLASS_GROUP_SEP + ginfo.Tag
	}
	return db.Elements[ginfo.Class].(*core.Class).Tag
}

// For testing
func printAtomicGroups(
	db *core.DbTopLevel,
	classDivs map[core.Ref][][]core.Ref,
	atomicGroups map[core.Ref][]AtomicGroup,
) {
	for _, cl := range db.Classes {
		agls := []string{}
		for _, ag := range atomicGroups[cl.Id] {
			agls = append(agls, ag.Tag)
		}
		fmt.Printf("  ++ %s: %+v\n", cl.Tag, agls)
		for _, div := range classDivs[cl.Id] {
			for _, gref := range div {
				agls := []string{}
				for _, ag := range atomicGroups[gref] {
					agls = append(agls, ag.Tag)
				}
				fmt.Printf("    -- %s: %+v\n", classOrGroup(db, gref), agls)
			}
		}
	}
}
