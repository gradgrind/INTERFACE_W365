package fet

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
