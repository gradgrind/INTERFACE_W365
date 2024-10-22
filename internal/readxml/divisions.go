package readxml

/* The only Groups relevant for timetabling are those listed under the
classes' Divisions ("GradePartiton" in Waldorf 365). The Name field specifies
a name for the Division.

Note that Waldorf 365 also uses Classes themselves as groups when specifying
the students taking part in a course.

TODO: The Divisions could be further "thinned out" by removing all those
whose Groups (all of them!) are not used by Courses?

TODO: The ordering of Groups within a division is not yet clear. Is it given
by the order within the Groups field of the Division element? Or must the
ListPosition of the Groups be used to sort them? Do these deliver the same
result?

TODO: There is also the question of the order of Divisions within a Class, but
this may not be very important.
*/

/*
import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type ClassDivision struct {
	Name   string
	Groups []string // list of Ids
}

func read_divisions(idmap IdMap, class_id string) []ClassDivision {
	divs := []ClassDivision{}
	class_node := idmap.Id2Node[class_id].(*Class)
	cl := fmt.Sprintf("%d%s", class_node.Level, class_node.Letter)
	divlist := []string{}
	if len(class_node.Divisions) != 0 {
		// First sort divisions
		divisions := []*Division{}
		for _, div_id := range strings.Split(class_node.Divisions, ",") {
			div_, ok := idmap.Id2Node[div_id]
			if !ok {
				log.Printf("Undefined Division in %s: %s\n", cl, div_id)
			} else {
				divisions = append(divisions, div_.(*Division))
			}
		}
		slices.SortFunc(divisions, func(a, b *Division) int {
			if a.ListPosition < b.ListPosition {
				return -1
			}
			if a.ListPosition == b.ListPosition {
				return 0
			}
			return 1
		})
		// Collect groups within each division
		for _, div := range divisions {
			glist := []string{}
			if len(div.Groups) == 0 {
				glist = append(glist, "*")
			} else {
				groupids := []string{}
				// First sort groups
				groups := []*Group{}
				for _, g_id := range strings.Split(div.Groups, ",") {
					g_, ok := idmap.Id2Node[g_id]
					if !ok {
						log.Printf("Undefined Group in %s: %s\n", cl, g_id)
					} else {
						g := g_.(*Group)
						//fmt.Printf("   ++LP%s.%s: %f\n", cl, g.Shortcut, g.ListPosition)
						groups = append(groups, g)
					}
				}
				slices.SortFunc(groups, func(a, b *Group) int {
					if a.ListPosition < b.ListPosition {
						return -1
					}
					if a.ListPosition == b.ListPosition {
						return 0
					}
					return 1
				})
				// Collect the groups
				for _, g := range groups {
					gstr, _ := groupTagFull(idmap, g.Id)
					glist = append(glist, gstr)
					groupids = append(groupids, g.Id)
				}
				if len(groupids) > 1 {
					divs = append(divs, ClassDivision{div.Name, groupids})
				} else {
					log.Printf("Incomplete Division in %s: %s (%d)\n",
						cl, div.Name, len(groupids))
				}
			}
			divlist = append(divlist, strings.Join(glist, ","))
		}
	}

	fmt.Printf("  Class %s: %s\n", cl, strings.Join(divlist, "//"))

	return divs
}
*/
