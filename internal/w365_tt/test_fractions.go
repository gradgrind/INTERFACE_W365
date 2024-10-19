package w365_tt

import (
	"fmt"
	"log"
	"strings"
)

func read_used_fractions(w365 *W365TT, idmap IdMap) map[string]int {
	fractions := map[string]int{}
	for _, l := range w365.Lessons {
		if len(l.Fractions) == 0 {
			continue
		}
		for _, x := range strings.Split(l.Fractions, ",") {
			f_, ok := idmap.Id2Node[x]
			if !ok {
				log.Printf("Lesson %s: No Fraction %s\n", l.Id, x)
				continue
			}
			f, ok := f_.(*Fraction)
			if !ok {
				log.Printf("Lesson %s: Not a Fraction %s\n", l.Id, x)
				continue
			}
			fractions[f.Id] += 1 // f.Id == x
		}
	}
	return fractions
}

func test_fractions(w365 *W365TT, idmap IdMap) {
	fractions := read_used_fractions(w365, idmap)
	for _, f := range w365.Fractions {
		fmt.Printf(" +++++++++++ Fraction %s:\n", f.Id)
		glist := []string{}
		if fractions[f.Id] == 0 {
			log.Fatalln("     ... is not used")
		}
		if len(f.SuperGroups) == 0 {
			log.Fatalln("     ... has no SuperGroups")
			continue
		}
		cl := ""
		for _, x := range strings.Split(f.SuperGroups, ",") {
			// Here a Class is not OK
			s_, ok := idmap.Id2Node[x]
			if !ok {
				log.Printf("     ... bad SuperGroup %s\n", x)
				continue
			}
			g, ok := s_.(*Group)
			if !ok {
				log.Printf("     ... not a Group %s\n", x)
				continue
			}
			cl1 := idmap.Group2Class[x].Tag()
			if cl1 != cl {
				if cl == "" {
					cl = cl1
				} else {
					cg, _ := groupTagFull(idmap, x)
					log.Fatalf("     ... Class mismatch (%s)\n", cg)
				}
			}
			glist = append(glist, g.Shortcut)
		}
		fmt.Printf(" === %s: %s\n", cl, strings.Join(glist, ","))
	}
}
