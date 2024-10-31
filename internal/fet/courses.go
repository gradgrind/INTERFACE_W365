package fet

import (
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
	"slices"
)

func gatherCourseGroups(fetinfo *fetInfo) {
	db := fetinfo.db
	fetinfo.superSubs = make(map[w365tt.Ref][]w365tt.Ref)
	fetinfo.courseGroups = make(map[w365tt.Ref][]w365tt.Ref)
	// Collect Courses with Lessons.
	for _, l := range db.Lessons {
		lcref := l.Course
		_, ok := fetinfo.courseGroups[lcref]
		if ok {
			continue
		}
		c := db.Elements[lcref] // can be Course or SuperCourse
		cnode, ok := c.(*w365tt.Course)
		if ok {
			fetinfo.courseGroups[lcref] = cnode.Groups
		} else {
			_, ok = c.(*w365tt.SuperCourse)
			if ok {
				fetinfo.superSubs[lcref] = []w365tt.Ref{}
				fetinfo.courseGroups[lcref] = []w365tt.Ref{}
				continue
			}
			log.Fatalf(
				"*ERROR* Invalid Course in Lesson %s:\n  %s\n",
				l.Id, lcref)
		}
	}
	// Now find the SubCourses
	for _, sbc := range db.SubCourses {
		spc := sbc.SuperCourse
		sblist, ok := fetinfo.superSubs[spc]
		if ok {
			// Only fill SuperCourses which have Lessons
			fetinfo.superSubs[spc] = append(sblist, sbc.Id)
			cglist := append(fetinfo.courseGroups[spc], sbc.Groups...)
			slices.Sort(cglist)
			cglist = slices.Compact(cglist)
			cglx := make([]w365tt.Ref, len(cglist))
			copy(cglx, cglist)
			fetinfo.courseGroups[spc] = cglx
		}
	}
}

// TODO-- Deprecated?
func readCourseIndexes(fetinfo *fetInfo) {
	courses := map[Ref]int{}
	for i, c := range fetinfo.db.Courses {
		courses[c.Id] = i
	}
	fetinfo.courses = courses
	supercourses := map[Ref]int{}
	for i, c := range fetinfo.db.SuperCourses {
		supercourses[c.Id] = i
	}
	fetinfo.supercourses = supercourses
	subcourses := map[Ref]int{}
	supersubs := map[Ref][]Ref{}
	for i, c := range fetinfo.db.SubCourses {
		subcourses[c.Id] = i
		cs := c.SuperCourse
		supersubs[cs] = append(supersubs[cs], c.Id)
	}
	fetinfo.subcourses = subcourses
	fetinfo.supersubs = supersubs
}
