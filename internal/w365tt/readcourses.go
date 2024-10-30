package w365tt

import (
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
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

func (dbp *DbTopLevel) newRoomChoice() string {
	// A rather primitive new-roomchoice-tag generator
	i := 0
	for {
		i++
		tag := "[" + strconv.Itoa(i) + "]"
		_, nok := dbp.RoomChoiceTags[tag]
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

func (dbp *DbTopLevel) readCourse(course *Course) {
	// Deal with subject
	//	var sr Ref = 0
	var ok bool
	msg1 := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
	msg2 := "*ERROR* Course %s:\n  Not a Subject: %s\n"
	if course.Subject == "" {
		if len(course.Subjects) == 1 {
			wsid := course.Subjects[0]
			s0, ok := dbp.Elements[wsid]
			if !ok {
				log.Fatalf(msg1, course.Id, wsid)
			}
			if _, ok = s0.(Subject); !ok {
				log.Fatalf(msg2, course.Id, wsid)
			}
		} else if len(course.Subjects) > 1 {
			// Make a subject name
			sklist := []string{}
			for _, wsid := range course.Subjects {
				// Need Tag/Shortcut field
				s0, ok := dbp.Elements[wsid]
				if ok {
					s, ok := s0.(Subject)
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
			_, ok = s0.(Subject)
			if !ok {
				log.Fatalf(msg2, course.Id, wsid)
			}
		} else {
			log.Fatalf(msg1, course.Id, wsid)
		}
	}

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
		_, ok = g.(Group)
		if !ok {
			// Check for class.
			_, ok = g.(Class)
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
		_, ok = t.(Teacher)
		if !ok {
			log.Fatalf("*ERROR* Invalid teacher in Course %s:\n  %s\n",
				course.Id, tref)
			//continue
		}
		//tlist = append(tlist, tref)
	}
	// Deal with rooms. W365 can have a single RoomGroup or a list of Rooms.
	rclist := []Ref{}     // choice list
	taglist := []string{} // list of room Tags/Shortcuts
	var rm Ref            // actual "room"
	for _, rref := range course.PreferredRooms {
		r, ok := dbp.Elements[rref]
		if !ok {
			log.Fatalf(
				"*ERROR* Unknown preferred room in Course %s:\n  %s\n",
				course.Id, rref)
			//continue
		}
		rr, ok := r.(Room)
		if ok {
			rclist = append(rclist, rref)
			taglist = append(taglist, rr.Tag)
		} else {
			_, ok = r.(RoomGroup)
			if !ok {
				log.Fatalf(
					"*ERROR* Invalid preferred room in Course %s:\n  %s\n",
					course.Id, rref)
				//continue
			}
			rclist = []Ref{rref}
			if len(course.PreferredRooms) != 1 {
				log.Printf(
					"*ERROR* Mixed Rooms and RoomGroups in Course %s\n",
					course.Id)
			}
			break
		}
	}
	if len(rclist) == 1 {
		// Take the single Room.
		rm = rclist[0]
	} else if len(rclist) > 1 {
		// Need a RoomChoiceGroup.
		// Reuse these if the same list appears again, but treat the room
		// order as significant.
		rs := strings.Join(taglist, ",")
		rm, ok = dbp.RoomChoiceTags[rs]
		if !ok {
			rk := fmt.Sprintf("RC%03d", len(dbdata.roomchoices)+1)
			rm = dbdata.nextId()
			dbp.RoomChoiceGroups = append(
				dbp.RoomChoiceGroups, db.RoomChoiceGroup{
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
