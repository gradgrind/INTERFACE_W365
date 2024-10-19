package w365_tt

import (
	"fmt"
	"log"
	"strings"
)

func test_courses(w365 *W365TT, idmap IdMap) map[string]int {
	used_groups := map[string]int{}
	for _, c := range w365.EpochPlanCourses {
		if len(c.Groups) == 0 {
			log.Printf("EpochPlanCourse with no groups: %s\n", c.Id)
			continue
		}
		for _, x := range strings.Split(c.Groups, ",") {
			g_, ok := idmap.Id2Node[x]
			if !ok {
				log.Printf("EpochPlanCourse %s: No Group %s\n", c.Id, x)
				continue
			}
			g, ok := g_.(*Group)
			if ok {
				used_groups[g.Id] += 1 // g.Id == x
				continue
			}
			cl, ok := g_.(*Class)
			if ok {
				fmt.Printf("  ++ Course Class %s\n", cl.Tag())
				continue
			}
			log.Printf("EpochPlanCourse %s: Group %s\n", c.Id, x)
		}
	}
	for _, c := range w365.Courses {
		if len(c.Groups) == 0 {
			log.Printf("Course with no groups: %s\n", c.Id)
			continue
		}
		for _, x := range strings.Split(c.Groups, ",") {
			g_, ok := idmap.Id2Node[x]
			if !ok {
				log.Printf("Course %s: No Group %s\n", c.Id, x)
				continue
			}
			g, ok := g_.(*Group)
			if ok {
				used_groups[g.Id] += 1 // g.Id == x
				continue
			}
			cl, ok := g_.(*Class)
			if ok {
				used_groups[cl.Id] += 1 // cl.Id == x
				//fmt.Printf("  ++ Course Class %s\n", cl.Tag())
				continue
			}
			log.Printf("Course %s: Group %s\n", c.Id, x)
		}
	}
	return used_groups
}
