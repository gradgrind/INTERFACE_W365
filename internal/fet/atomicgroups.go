package fet

import (
	"fmt"
	"log"
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

func filterDivisions(fetinfo *fetInfo) {
	// Prepare filtered versions of the class Divisions containing only
	// those Divisions which have Groups used in Lessons.

	// Collect groups used in Lessons. Get them from the
	// fetinfo.courseGroups map, which only includes courses with lessons.
	usedgroups := map[Ref]bool{}
	for _, cg := range fetinfo.courseGroups {
		for _, g := range cg {
			usedgroups[g] = true
		}
	}
	// Filter the class divisions, discarding the division names.
	cdivs := map[Ref][][]Ref{}
	for _, c := range fetinfo.db.Classes {
		divs := [][]Ref{}
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
	fetinfo.classDivisions = cdivs
}

func makeAtomicGroups(fetinfo *fetInfo) {
	// An atomic group is an ordered list of single groups from each division.
	fetinfo.atomicgroups = map[Ref][]AtomicGroup{}
	// Go through the classes inspecting their Divisions. Retain only those
	// which have lessons.
	filterDivisions(fetinfo) // -> fetinfo.classDivisions
	// Go through the classes inspecting their Divisions.
	// Build a list-basis for the atomic groups based on the Cartesian product.
	for _, cl := range fetinfo.db.Classes {
		divs, ok := fetinfo.classDivisions[cl.Id]
		if !ok {
			log.Fatalf("*BUG* fetinfo.classDivisions[%s]\n", cl.Id)
		}
		// The atomic groups will be built as a list of lists of Refs.
		agrefs := [][]Ref{{}}
		for _, dglist := range divs {
			// Add another division – increases underlying list lengths.
			agrefsx := [][]Ref{}
			for _, ag := range agrefs {
				// Extend each of the old list items by appending each
				// group of the new division in turn – multiplies the
				// total number of atomic groups.
				for _, g := range dglist {
					gx := make([]Ref, len(ag)+1)
					copy(gx, append(ag, g))
					agrefsx = append(agrefsx, gx)
				}
			}
			agrefs = agrefsx
		}
		//fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		//fmt.Printf("     --> %+v\n", agrefs)

		for _, cl := range fetinfo.db.Classes {
			divs, ok := fetinfo.classDivisions[cl.Id]
			if !ok {
				log.Fatalf("*ERROR* fetinfo.classDivisions[%s]\n", cl.Id)
			}
			// Make AtomicGroups
			aglist := []AtomicGroup{}
			for _, ag := range agrefs {
				glist := []string{}
				for _, gref := range ag {
					glist = append(glist, fetinfo.ref2grouponly[gref])
				}
				ago := AtomicGroup{
					Class:  cl.Id,
					Groups: ag,
					Tag: fmt.Sprintf(
						"%s#%s", cl.Tag, strings.Join(glist, "~")),
				}
				aglist = append(aglist, ago)
			}

			fmt.Printf("   %s ++> %+v\n", cl.Tag, aglist)

			// Map the individual groups to their atomic groups.
			g2ags := map[Ref][]AtomicGroup{}
			//		xg2ags := map[string][]string{}
			i := len(divs)
			n := 1
			for i > 0 {
				i--
				a := 0 // ag index

				for a < len(aglist) {
					for x, g := range divs[i] {
						for j := 0; j < n; j++ {

							fmt.Printf(" ????? cl=%s, a=%d, i=%d, x=%d, j=%d\n",
								cl.Tag, a, i, x, j)
							fmt.Printf("   -- agl=%d, xl=%d, n=%d\n",
								len(aglist), len(divs), n)

							g2ags[g] = append(g2ags[g], aglist[a])
							//	g2ags[fetinfo.ref2fet[g]] = append(
							//	    xg2ags[fetinfo.ref2fet[g]], aglist[a])
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
}

func printAtomicGroups(fetinfo *fetInfo) {
	for _, cl := range fetinfo.db.Classes {
		agls := []string{}
		for _, ag := range fetinfo.atomicgroups[cl.Id] {
			agls = append(agls, ag.Tag)
		}
		fmt.Printf("  ++ %s: %+v\n", fetinfo.ref2fet[cl.Id], agls)
		for _, div := range fetinfo.classDivisions[cl.Id] {
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
