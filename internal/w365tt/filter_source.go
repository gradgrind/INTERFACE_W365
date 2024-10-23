package w365tt

import (
	"fmt"
	"log"
	"strings"
)

// Courses with more than one subject:
// Create a new subject (e.g. "Xn"). It's Name field could be built from
// the list of subjects. Use a cache to check for repeats. Replace the lists
// by the new subject.
// This is a preliminary implementation of the function. In the end a
// somewhat different approach might be necessary.
func Multisubjects(w365 *W365TopLevel) {
	// First gather references to all Subject nodes.
	//TODO: Is this already available somewhere?
	id2s := map[W365Ref]Subject{}
	for _, s := range w365.Subjects {
		id2s[s.IdStr()] = s
	}
	cache := map[string]W365Ref{}
	// Now check all Courses and SubCourses for multiple subjects.
	n := 0
	for i, c := range w365.Courses {
		if len(c.Subjects) <= 1 {
			continue
		}
		var slist []string
		for _, s := range c.Subjects {
			snode, ok := id2s[s]
			if !ok {
				log.Fatalf("Course %s has invalid Subject %s\n", c.Id, s)
			}
			slist = append(slist, snode.Shortcut)
		}
		sname := strings.Join(slist, ",")
		sid, ok := cache[sname]
		if !ok {
			n++
			tag := fmt.Sprintf("X%02d", n)
			sid = W365Ref(fmt.Sprintf("Id_%s", tag))
			w365.Subjects = append(w365.Subjects, Subject{
				Id:       sid,
				Type:     TypeSUBJECT,
				Name:     sname,
				Shortcut: tag,
			})
			cache[sname] = sid
		}
		w365.Courses[i].Subjects = []W365Ref{sid}
	}
	for i, c := range w365.SubCourses {
		if len(c.Subjects) <= 1 {
			continue
		}
		var slist []string
		for _, s := range c.Subjects {
			snode, ok := id2s[s]
			if !ok {
				log.Fatalf("SubCourse %s has invalid Subject %s\n", c.Id, s)
			}
			slist = append(slist, snode.Shortcut)
		}
		sname := strings.Join(slist, ",")
		sid, ok := cache[sname]
		if !ok {
			n++
			tag := fmt.Sprintf("X%02d", n)
			sid = W365Ref(fmt.Sprintf("Id_%s", tag))
			w365.Subjects = append(w365.Subjects, Subject{
				Id:       sid,
				Type:     TypeSUBJECT,
				Name:     sname,
				Shortcut: tag,
			})
			cache[sname] = sid
		}
		w365.SubCourses[i].Subjects = []W365Ref{sid}
	}
}

// Might not be necessary: Filter out all Divisions where none of the groups
// are used by a course.

// This version doesn't actually strip, it just reports.
func StripDivisions(w365 *W365TopLevel) {
	// First gather references to all Group nodes, setting their count to 0.
	gcount := map[W365Ref]int{}
	for _, g := range w365.Groups {
		gcount[g.IdStr()] = 0
	}
	// Now collect the used Groups (not Classes)
	for _, crs := range w365.Courses {
		for _, gr := range crs.Groups {
			if _, ok := gcount[gr]; ok {
				gcount[gr]++
			}
		}
	}
	// Now go through the classes checking the divisions.
	for _, c := range w365.Classes {
		for _, d := range c.Divisions {
			used := false
			for i, gr := range d.Groups {
				if gcount[gr] != 0 {
					used = true
					fmt.Printf("** %s / %s, %d: %d\n",
						c.Shortcut, d.Name, i+1, gcount[gr])
				} else {
					log.Printf("In Class %s, Division %s:\n  Group %s not used.\n",
						c.Shortcut, d.Name, gr)
				}
			}
			if !used {
				log.Printf("In Class %s: Division %s has no courses.\n",
					c.Shortcut, d.Name)
			}
		}
	}
}
