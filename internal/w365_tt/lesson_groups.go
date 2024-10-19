package w365_tt

import (
	"fmt"
	"log"
	"strings"
)

func test_lesson_groups(w365 *W365TT, idmap IdMap) {
	for _, l := range w365.Lessons {
		if len(l.Course) == 0 {
			// Let's assume this Lesson is invalid/unused
			log.Fatalf("  !!! Lessons.NoCourse %+v\n", l)
			continue
		} else {
			crs, ok := idmap.Id2Node[l.Course]
			if !ok {
				// Let's assume this Lesson is invalid/unused
				log.Printf("  !!! Lessons.Course %+v\n", l)
				continue
			} else {
				var groups []string
				crsc, ok := crs.(*Course)
				if ok {
					// Get the Groups
					if len(crsc.Groups) == 0 {
						log.Printf("Course without Groups: %s\n", crsc.Id)
						continue
					}
					for _, g := range strings.Split(crsc.Groups, ",") {
						gstr, _ := groupTagFull(idmap, g)
						groups = append(groups, gstr)
					}
					fmt.Printf("***** Lesson %s\n  ++ Course %s: %s\n",
						l.Id, crsc.Id, strings.Join(groups, ","))
				} else {
					crsc, ok := crs.(*EpochPlanCourse)
					if ok {
						// Get the Groups
						if len(crsc.Groups) == 0 {
							log.Printf("Course without Groups: %s\n", crsc.Id)
							continue
						}
						for _, g := range strings.Split(crsc.Groups, ",") {
							gstr, _ := groupTagFull(idmap, g)
							groups = append(groups, gstr)
						}
						fmt.Printf("***** Lesson %s\n  ++ EpochPlanCourse %s: %s\n",
							l.Id, crsc.Id, strings.Join(groups, ","))
					} else {
						log.Printf("  !!! Lessons.Course? %s\n", l.Course)
					}
				}
				// Now show the Fractions
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
					if len(f.SuperGroups) == 0 {
						log.Fatalln("     ... has no SuperGroups")
						continue
					}
					var glist []string
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
		}
	}
}
