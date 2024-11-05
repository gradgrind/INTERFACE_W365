package readxml

import (
	"gradgrind/INTERFACE_W365/internal/w365tt"
	"log"
)

//TODO: The Course field can be a Course or a SuperCourse.
// Only lessons placed in a schedule are present, and only one schedule
// should be used. Perhaps regard all schedulued lessons as fixed regardless
// of the flag?
// Non-scheduled lessons must be discovered (from the Courses and SuperCourses)
// and then added.
// Bear in mind that lessons have no length field, so multi-hour lessons are
// made up of single lessons. It might be worth replacing these by real
// multi-hour lessons? Of course, they must be reverted to W365 form to pass
// them back , but this should be possible ...

func readLessons(
	//outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	items []Lesson,
	scheduled []w365tt.Ref,
) map[w365tt.Ref][]*Lesson {
	courseLlist := map[w365tt.Ref][]*Lesson{} // course ref -> lesson list
	for _, n := range items {
		nid := addId(id2node, &n)
		if nid == "" {
			continue
		}
	}
	for _, sl := range scheduled {
		n, ok := id2node[sl]
		if !ok {
			log.Printf("*ERROR* Lesson in Schedule has no Definition: %s\n",
				sl)
			continue
		}
		np, ok := n.(*Lesson)
		if !ok {
			log.Printf("*ERROR* Bad Lesson in Schedule: %s\n", sl)
			continue
		}
		if _, ok := id2node[np.Course]; !ok {
			log.Printf("*ERROR* Lesson with invalid Course: %s\n", np.Id)
			continue
		}
		//msg := fmt.Sprintf("Course %s in LocalRooms", nid)
		//rlist := GetRefList(id2node, n.LocalRooms, msg)
		courseLlist[np.Course] = append(courseLlist[np.Course], np)

		/*
			outdata.Lessons = append(outdata.Lessons, w365tt.Lesson{
						Id:       nid,
						Course:   n.Course,
						Duration: dur,
						Day:      n.Day,
						Hour:     n.Hour,
						Fixed:    n.Fixed,
						Rooms:    GetRefList(id2node, n.LocalRooms, msg),
					})
		*/
	}
	return courseLlist
}

// TODO
func makeLessons(
	outdata *w365tt.DbTopLevel,
	id2node map[w365tt.Ref]interface{},
	courseLessons map[w365tt.Ref][]int,
	lessons map[w365tt.Ref][]*Lesson,
) {

}
