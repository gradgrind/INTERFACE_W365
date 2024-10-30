package w365tt

import (
	"log"
	"strconv"
	"strings"
)

func (dbp *DbTopLevel) readSubjects() {
	for _, n := range dbp.Subjects {
		_, nok := dbp.SubjectTags[n.Tag]
		if nok {
			log.Fatalf("*ERROR* Subject Tag (Shortcut) defined twice: %s\n",
				n.Tag)
		}
		t, nok := dbp.SubjectNames[n.Name]
		if nok {
			log.Printf("*WARNING* Subject Name defined twice (different"+
				" Tag/Shortcut):\n  %s (%s/%s)\n", n.Name, t, n.Tag)
		} else {
			dbp.SubjectNames[n.Name] = n.Tag
		}
		dbp.SubjectTags[n.Tag] = n.Id
	}
}

func (dbp *DbTopLevel) newSubject() string {
	// A rather primitive new-subject-tag generator
	i := 0
	for {
		i++
		tag := "X" + strconv.Itoa(i)
		_, nok := dbp.SubjectTags[tag]
		if !nok {
			return tag
		}
	}
}

func (dbp *DbTopLevel) readCourses() {
	for i := 0; i < len(dbp.Courses); i++ {
		n := &dbp.Courses[i]
		dbp.Elements[n.Id] = n

		dbp.readCourse(n)
	}
}

func (dbp *DbTopLevel) checkCourseSubject(course *Course) {
	// Deal with the subject(s) fields
	//	var sr Ref = 0
	msg1 := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
	msg2 := "*ERROR* Course %s:\n  Not a Subject: %s\n"
	if course.Subject == "" {
		if len(course.Subjects) == 1 {
			wsid := course.Subjects[0]
			s0, ok := dbp.Elements[wsid]
			if !ok {
				log.Fatalf(msg1, course.Id, wsid)
			}
			if _, ok = s0.(*Subject); !ok {
				log.Fatalf(msg2, course.Id, wsid)
			}
		} else if len(course.Subjects) > 1 {
			// Make a subject name
			sklist := []string{}
			for _, wsid := range course.Subjects {
				// Need Tag/Shortcut field
				s0, ok := dbp.Elements[wsid]
				if ok {
					s, ok := s0.(*Subject)
					if !ok {
						log.Fatalf(msg2, course.Id, wsid)
					}
					sklist = append(sklist, s.Tag)
				} else {
					log.Fatalf(msg1, course.Id, wsid)
				}
			}
			skname := strings.Join(sklist, ",")
			stag, ok := dbp.SubjectNames[skname]
			if ok {
				// The Name has already been used.
				course.Subject = dbp.SubjectTags[stag]
			} else {
				// Need a new Subject.
				stag = dbp.newSubject()
				sref := dbp.NewId()
				i := len(dbp.Subjects)
				dbp.Subjects = append(dbp.Subjects, Subject{
					Id:   sref,
					Tag:  stag,
					Name: skname,
				})
				dbp.AddElement(sref, &dbp.Subjects[i])
				dbp.SubjectTags[stag] = sref
				dbp.SubjectNames[skname] = stag
				course.Subject = sref
			}
			// Clear Subjects field.
			course.Subjects = nil
		}
	} else {
		if len(course.Subjects) != 0 {
			log.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
				course.Id)
		}
		wsid := course.Subject
		s0, ok := dbp.Elements[wsid]
		if ok {
			_, ok = s0.(*Subject)
			if !ok {
				log.Fatalf(msg2, course.Id, wsid)
			}
		} else {
			log.Fatalf(msg1, course.Id, wsid)
		}
	}
}

func (dbp *DbTopLevel) readCourse(course *Course) {
	dbp.checkCourseSubject(course)
	// Deal with groups
	//glist := []Ref{}
	for _, gref := range course.Groups {
		g, ok := dbp.Elements[gref]
		if !ok {
			log.Fatalf("*ERROR* Unknown group in Course %s:\n  %s\n",
				course.Id, gref)
			//continue
		}
		// g can be a Group or a Class.
		_, ok = g.(*Group)
		if !ok {
			// Check for class.
			_, ok = g.(*Class)
			if !ok {
				log.Fatalf("*ERROR* Invalid group in Course %s:\n  %s\n",
					course.Id, gref)
				//continue
			}
		}
		//glist = append(glist, gref)
	}
	// Deal with teachers
	//tlist := []Ref{}
	for _, tref := range course.Teachers {
		t, ok := dbp.Elements[tref]
		if !ok {
			log.Fatalf("*ERROR* Unknown teacher in Course %s:\n  %s\n",
				course.Id, tref)
			//continue
		}
		_, ok = t.(*Teacher)
		if !ok {
			log.Fatalf("*ERROR* Invalid teacher in Course %s:\n  %s\n",
				course.Id, tref)
			//continue
		}
		//tlist = append(tlist, tref)
	}
	// Deal with rooms. W365 can have a single RoomGroup or a list of Rooms.

	rref := Ref("")
	if len(course.PreferredRooms) > 1 {
		// Make a RoomChoiceGroup
		var estr string
		rref, estr = dbp.makeRoomChoiceGroup(course.PreferredRooms)
		if estr != "" {
			log.Printf("*ERROR* In Course %s:\n%s", course.Id, estr)
		}
	} else if len(course.PreferredRooms) == 1 {
		// Check that room is Room or RoomGroup.
		rref0 := course.PreferredRooms[0]
		r, ok := dbp.Elements[rref0]
		if ok {
			_, ok = r.(*Room)
			if ok {
				rref = rref0
			} else {
				_, ok = r.(*RoomGroup)
				if ok {
					rref = rref0
				} else {
					log.Printf("*ERROR* Invalid room in Course %s:\n  %s\n",
						course.Id, rref0)
				}
			}
		} else {
			log.Printf("*ERROR* Unknown room in Course %s:\n  %s\n",
				course.Id, rref0)
		}
	}
	if course.Room != "" {
		if rref != "" {
			log.Printf(
				"*ERROR* Course has both Room and Rooms entries:\n %s\n",
				course.Id)
		}
		r, ok := dbp.Elements[course.Room]
		if ok {
			_, ok = r.(*Room)
			if !ok {
				_, ok = r.(*RoomGroup)
				if !ok {
					_, ok = r.(*RoomChoiceGroup)
					if !ok {
						log.Printf(
							"*ERROR* Invalid room in Course %s:\n  %s\n",
							course.Id, course.Room)
						course.Room = ""
					}
				}
			}

		} else {
			log.Printf("*ERROR* Unknown room in Course %s:\n  %s\n",
				course.Id, course.Room)
			course.Room = ""
		}
	} else {
		course.Room = rref
	}
	course.PreferredRooms = nil
}
