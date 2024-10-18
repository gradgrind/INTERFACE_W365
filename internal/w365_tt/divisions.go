package w365_tt

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

func read_divisions(idmap IdMap, class_id string) {
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
		for _, div := range divisions {
			glist := []string{}
			if len(div.Groups) == 0 {
				glist = append(glist, "*")
			} else {
				// First sort groups
				groups := []*Group{}
				for _, g_id := range strings.Split(div.Groups, ",") {
					g_, ok := idmap.Id2Node[g_id]
					if !ok {
						log.Printf("Undefined Group in %s: %s\n", cl, g_id)
					} else {
						g := g_.(*Group)

						fmt.Printf("   ++LP%s.%s: %f\n", cl, g.Shortcut, g.ListPosition)

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
				for _, g := range groups {
					gstr, _ := groupTagFull(idmap, g.Id)
					glist = append(glist, gstr)
				}
			}
			divlist = append(divlist, strings.Join(glist, ","))
		}
	}

	fmt.Printf("  %s: %s\n", cl, strings.Join(divlist, "//"))
}
