package w365tt

import (
	"log"
)

func (dbdata *xData) readSubjects() {
	for _, n := range dbdata.data.Subjects {
		_, nok := dbdata.subjecttags[n.Tag]
		if nok {
			log.Fatalf("*ERROR* Subject Tag (Shortcut) defined twice: %s\n",
				n.Tag)
		}
		t, nok := dbdata.subjectnames[n.Name]
		if nok {
			log.Printf("*WARNING* Subject Name defined twice (different"+
				" Tag/Shortcut):\n  %s (%s/%s)\n", n.Name, t, n.Tag)
		} else {
			dbdata.subjectnames[n.Name] = n.Tag
		}
		dbdata.subjecttags[n.Tag] = n.Id
	}
}

/*

func (dbdata *xData) readCourses() {
	for i := 0; i < len(dbdata.data.Courses); i++ {
		n := &dbdata.data.Courses[i]
		dbdata.elements[n.Id] = n

		dbdata.readCourse(n)
	}
}

func (dbdata *xData) readCourse(course *Course) {
	// Deal with subject
	//	var sr Ref = 0
	var ok bool
	msg := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
	msg2 := "*ERROR* Course %s:\n  Not a Subject: %s\n"
	if course.Subject == "" {
		if len(course.Subjects) == 1 {
			wsid := course.Subjects[0]
			s0, ok := dbdata.elements[wsid]
			if !ok {
				log.Fatalf(msg, course.Id, wsid)
			}
			if _, ok = s0.(Subject); !ok {
				log.Fatalf(msg2, course.Id, wsid)
			}
		} else if len(course.Subjects) > 1 {
			// Make a subject name
			sklist := []string{}
			for _, wsid := range course.Subjects {
				// Need Shortcut field
				s0, ok := dbdata.elements[wsid]
				if ok {
					s, ok := s0.(Subject)
					if !ok {
						log.Fatalf(msg2, course.Id, wsid)
					}
					sklist = append(sklist, s.Tag)
				} else {
					log.Fatalf(msg, course.Id, wsid)
				}
			}
			skname := strings.Join(sklist, ",")
			stag, ok := dbdata.subjectnames[skname]
            if ok {
                // The Name has already been used.
                course.Subject = dbdata.subjecttags[stag]
            } else {
                // Need a new Subject.

                ref = new
				s := Subject{
					Id:   sref,
					Tag:  sk,
					Name: skname,
				}

				sk := fmt.Sprintf("X%02d", len(dbdata.newsubjects)+1)
				sr = dbdata.nextId()
				dbdata.data.Subjects = append(dbdata.data.Subjects,
					db.Subject{
						Id:   sr,
						Tag:  sk,
						Name: skname,
					})
				dbdata.newsubjects[skname] = sr
			}
		}
	} else {
		if len(subjects) != 0 {
			log.Printf("*ERROR* Course has both Subject AND Subjects: %s\n", id)
		}
		wsid := subject
		sr, ok = dbdata.subjects[wsid]
		if !ok {
			log.Printf(msg, id, wsid)
		}
	}
	// Deal with groups
	glist := []db.DbRef{}
	for _, g := range groups {
		gr, ok := dbdata.groups[g]
		// gr can refer to a Group or a Class.
		if !ok {
			// Check for class.
			gr, ok = dbdata.classes[g]
			if !ok {
				log.Printf("*ERROR* Unknown group in Course %s:\n  %s\n", id, g)
				continue
			}
		}
		glist = append(glist, gr)
	}
	// Deal with teachers
	tlist := []db.DbRef{}
	for _, t := range teachers {
		tr, ok := dbdata.teachers[t]
		if !ok {
			log.Printf("*ERROR* Unknown teacher in Course %s:\n  %s\n", id, t)
			continue
		}
		tlist = append(tlist, tr)
	}
	// Deal with rooms. W365 can have a single RoomGroup or a list of Rooms
	rclist := []db.DbRef{} // choice list
	var rm db.DbRef        // actual "room"
	for _, r := range rooms {
		rr, ok := dbdata.rooms[r]
		if ok {
			rclist = append(rclist, rr)
		} else {
			rm, ok = dbdata.roomgroups[r]
			if ok {
				if len(rooms) != 1 {
					rclist = []db.DbRef{}
					log.Printf(
						"*ERROR* Mixed Rooms and RoomGroups in Course %s\n", id)
				}
				break
			} else {
				log.Printf("*ERROR* Unknown room in Course %s:\n  %s\n", id, r)
				continue
			}
		}
	}
	if len(rclist) == 1 {
		// Take the single Room.
		rm = rclist[0]
	} else if len(rclist) > 1 {
		// Need a RoomChoiceGroup.
		// Reuse these if the same list appears again, but treat the
		// order as significant.
		rslist := []string{}
		for _, r := range rclist {
			rslist = append(rslist, dbdata.roomtag[r])
		}
		rs := strings.Join(rslist, ",")
		rm, ok = dbdata.roomchoices[rs]
		if !ok {
			rk := fmt.Sprintf("RC%03d", len(dbdata.roomchoices)+1)
			rm = dbdata.nextId()
			dbdata.data.RoomChoiceGroups = append(
				dbdata.data.RoomChoiceGroups, db.RoomChoiceGroup{
					Id:    rm,
					Tag:   rk,
					Name:  rs,
					Rooms: rclist,
				})
			dbdata.roomchoices[rs] = rm
		}
	}
	return sr, glist, tlist, rm
}

*/
