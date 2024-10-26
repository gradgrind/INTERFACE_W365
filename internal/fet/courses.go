package fet

import "gradgrind/INTERFACE_W365/internal/db"

func readCourseIndexes(fetinfo *fetInfo) {
	courses := map[db.DbRef]int{}
	for i, c := range fetinfo.db.Courses {
		courses[c.Id] = i
	}
	fetinfo.courses = courses
	supercourses := map[db.DbRef]int{}
	for i, c := range fetinfo.db.SuperCourses {
		supercourses[c.Id] = i
	}
	fetinfo.supercourses = supercourses
	subcourses := map[db.DbRef]int{}
	supersubs := map[db.DbRef][]db.DbRef{}
	for i, c := range fetinfo.db.SubCourses {
		subcourses[c.Id] = i
		cs := c.SuperCourse
		supersubs[cs] = append(supersubs[cs], c.Id)
	}
	fetinfo.subcourses = subcourses
	fetinfo.supersubs = supersubs
}
