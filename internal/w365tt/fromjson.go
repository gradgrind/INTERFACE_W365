package w365tt

import (
	"encoding/json"
	"fmt"
	"gradgrind/INTERFACE_W365/internal/db"
	"io"
	"log"
	"os"
	"strings"
)

func ReadJSON(jsonpath string) W365TopLevel {
	// Open the  JSON file
	jsonFile, err := os.Open(jsonpath)
	if err != nil {
		log.Fatal(err)
	}
	// Remember to close the file at the end of the function
	defer jsonFile.Close()
	// read the opened XML file as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	log.Printf("*+ Reading: %s\n", jsonpath)
	v := W365TopLevel{}
	err = json.Unmarshal(byteValue, &v)
	if err != nil {
		log.Fatalf("Could not unmarshal json: %s\n", err)
	}

	DeMultipleSubjects(&v)

	return v
}

func defaultMinus1(v interface{}) int {
	if v == nil {
		return -1
	}
	return int(v.(float64))
}

type xData struct {
	w365        W365TopLevel
	data        db.DbTopLevel
	dbi         db.DbRef
	idmap       map[W365Ref]db.DbRef
	divgroups   map[db.DbRef]int    // count usage in courses
	classes     map[db.DbRef]int    // count usage in courses
	subjectmap  map[W365Ref]string  // Subject Tag (Shortcut)
	newsubjects map[string]db.DbRef // Subject's db Id
}

func (dbdata *xData) getId(wref W365Ref, msg string) db.DbRef {
	id, ok := dbdata.idmap[wref]
	if ok {
		return id
	}
	fmt.Printf("*ERROR* Unknown W365-Id: %s\n  ++ %s\n", wref, msg)
	return 0
}

func LoadJSON(jsonpath string) db.DbTopLevel {
	dbdata := xData{
		w365:  ReadJSON(jsonpath),
		data:  db.DbTopLevel{},
		dbi:   0,
		idmap: map[W365Ref]db.DbRef{},
	}

	dbdata.addInfo()
	dbdata.addDays()
	dbdata.addHours()
	dbdata.addTeachers()
	dbdata.addSubjects()
	dbdata.addRooms()
	// RoomChoicesGroups: W365 has none of these – they must be generated
	// from the PreferredRooms lists of courses.
	dbdata.data.RoomChoiceGroups = []db.RoomChoiceGroup{}
	// The RoomGroups from W365 are used by some courses. The listed rooms
	// should build the Rooms list for the course.
	dbdata.addGroups()
	dbdata.addClasses()

	return dbdata.data
}

func (dbdata *xData) nextId(w365Id W365Ref) db.DbRef {
	dbdata.dbi++
	if w365Id != "" {
		dbdata.idmap[w365Id] = dbdata.dbi
	}
	return dbdata.dbi
}

func (dbdata *xData) addInfo() {
	dbdata.data.Info = db.Info{
		FirstAfternoonHour: dbdata.w365.W365TT.FirstAfternoonHour,
		MiddayBreak:        dbdata.w365.W365TT.MiddayBreak,
	}
}

func (dbdata *xData) addDays() {
	for _, d := range dbdata.w365.Days {
		dbdata.data.Days = append(dbdata.data.Days, db.Day{
			Id:   dbdata.nextId(d.Id),
			Tag:  d.Shortcut,
			Name: d.Name,
		})
	}
}

func (dbdata *xData) addHours() {
	mdbok := len(dbdata.data.Info.MiddayBreak) == 0
	for i, d := range dbdata.w365.Hours {
		if d.FirstAfternoonHour {
			dbdata.data.Info.FirstAfternoonHour = i
		}
		if d.MiddayBreak {
			if mdbok {
				dbdata.data.Info.MiddayBreak = append(
					dbdata.data.Info.MiddayBreak, i)
			} else {
				fmt.Printf("*ERROR* MiddayBreak set in Info AND Hours\n")
			}
		}
		dbdata.data.Hours = append(dbdata.data.Hours, db.Hour{
			Id:    dbdata.nextId(d.Id),
			Tag:   d.Shortcut,
			Name:  d.Name,
			Start: d.Start,
			End:   d.End,
		})
	}
}

func (dbdata *xData) addTeachers() {
	for _, d := range dbdata.w365.Teachers {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		dbdata.data.Teachers = append(dbdata.data.Teachers, db.Teacher{
			Id:               dbdata.nextId(d.Id),
			Tag:              d.Shortcut,
			Name:             d.Name,
			Firstname:        d.Firstname,
			NotAvailable:     a,
			MinLessonsPerDay: defaultMinus1(d.MinLessonsPerDay),
			MaxLessonsPerDay: defaultMinus1(d.MaxLessonsPerDay),
			MaxDays:          defaultMinus1(d.MaxDays),
			MaxGapsPerDay:    defaultMinus1(d.MaxGapsPerDay),
			MaxGapsPerWeek:   defaultMinus1(d.MaxGapsPerWeek),
			MaxAfternoons:    defaultMinus1(d.MaxAfternoons),
			LunchBreak:       d.LunchBreak,
		})
	}
}

func (dbdata *xData) addSubjects() {
	dbdata.subjectmap = map[W365Ref]string{}
	for _, d := range dbdata.w365.Subjects {
		dbdata.data.Subjects = append(dbdata.data.Subjects, db.Subject{
			Id:   dbdata.nextId(d.Id),
			Tag:  d.Shortcut,
			Name: d.Name,
		})
		dbdata.subjectmap[d.Id] = d.Shortcut
	}
}

func (dbdata *xData) addRooms() {
	for _, d := range dbdata.w365.Rooms {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		dbdata.data.Rooms = append(dbdata.data.Rooms, db.Room{
			Id:           dbdata.nextId(d.Id),
			Tag:          d.Shortcut,
			Name:         d.Name,
			NotAvailable: a,
		})
	}
}

func (dbdata *xData) addGroups() {
	// Every Group must be a member of a Class Division.
	for _, d := range dbdata.w365.Groups {
		dbdata.data.Groups = append(dbdata.data.Groups, db.Group{
			Id:  dbdata.nextId(d.Id),
			Tag: d.Shortcut,
		})
	}
}

func (dbdata *xData) addClasses() {
	dbdata.divgroups = map[db.DbRef]int{}
	dbdata.classes = map[db.DbRef]int{}
	for _, d := range dbdata.w365.Classes {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		// Get the divisions and add their groups to divgroups, so that their
		// use in courses can be counted – for checking and filtering.
		divs := []db.Division{}
		for _, wdiv := range d.Divisions {
			glist := []db.DbRef{}
			for _, g := range wdiv.Groups {
				msg := fmt.Sprintf("Unknown Group in Class %s, Division %s",
					d.Shortcut, wdiv.Name)
				gr := dbdata.getId(g, msg)
				if gr != 0 {
					glist = append(glist, gr)
					dbdata.divgroups[gr] = 0
				}
			}
			divs = append(divs, db.Division{
				Name:   wdiv.Name,
				Groups: glist,
			})
		}
		cid := dbdata.nextId(d.Id)
		dbdata.data.Classes = append(dbdata.data.Classes, db.Class{
			Id:               cid,
			Name:             d.Name,
			Tag:              d.Shortcut,
			Level:            d.Level,
			Letter:           d.Letter,
			NotAvailable:     a,
			Divisions:        divs,
			MinLessonsPerDay: defaultMinus1(d.MinLessonsPerDay),
			MaxLessonsPerDay: defaultMinus1(d.MaxLessonsPerDay),
			MaxGapsPerDay:    defaultMinus1(d.MaxGapsPerDay),
			MaxGapsPerWeek:   defaultMinus1(d.MaxGapsPerWeek),
			MaxAfternoons:    defaultMinus1(d.MaxAfternoons),
			LunchBreak:       d.LunchBreak,
			ForceFirstHour:   d.ForceFirstHour,
		})
		dbdata.classes[cid] = 0
	}
}

func (dbdata *xData) addCourses() {
	for _, d := range dbdata.w365.Courses {
		// Deal with subject
		var sr db.DbRef = 0
		msg := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
		if d.Subject == "" {
			if len(d.Subjects) == 1 {
				wsid := d.Subjects[0]
				if _, ok := dbdata.subjectmap[wsid]; ok {
					sr = dbdata.idmap[wsid]
				} else {
					fmt.Printf(msg, d.Id, wsid)
				}
			} else if len(d.Subjects) > 1 {
				// Make a subject name
				sklist := []string{}
				for _, wsid := range d.Subjects {
					// Need Shortcut field
					sk, ok := dbdata.subjectmap[wsid]
					if ok {
						sklist = append(sklist, sk)
					} else {
						fmt.Printf(msg, d.Id, wsid)
					}
				}
				skname := strings.Join(sklist, ",")
				sid, ok := dbdata.newsubjects[skname]
				if !ok {
					sk := fmt.Sprintf("X%02d", len(dbdata.newsubjects)+1)
					sid = dbdata.nextId("")
					dbdata.data.Subjects = append(dbdata.data.Subjects,
						db.Subject{
							Id:   sid,
							Tag:  sk,
							Name: skname,
						})
					dbdata.newsubjects[skname] = sid
				}
			}
		} else {
			if len(d.Subjects) != 0 {
				fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
					d.Id)
			}
			wsid := d.Subject
			if _, ok := dbdata.subjectmap[wsid]; ok {
				sr = dbdata.idmap[wsid]
			} else {
				fmt.Printf(msg, d.Id, wsid)
			}
		}
		// Deal with groups
		//TODO: Only increment the counters if the course has lessons!
		glist := []db.DbRef{}
		for _, g := range d.Groups {
			msg := fmt.Sprintf("*ERROR* Unknown group in Course %s", d.Id)
			gid := dbdata.getId(g, msg)
			if gid != 0 {
				if _, ok := dbdata.divgroups[gid]; ok {
					dbdata.divgroups[gid]++
					//TODO: Could it be a group, but not in a division?
					// If so, that would be ok, but the course shouldn't
					// have any lessons!
				} else if _, ok = dbdata.classes[gid]; ok {
					dbdata.classes[gid]++
				} else {
					fmt.Printf("*ERROR* In Course %s,\n"+
						"  -- Element is not a valid Group/Class: %s", d.Id, g)
					continue
				}
				glist = append(glist, gid)
			}
		}
		// Deal with teachers
		tlist := []db.DbRef{}
		/*
			        for _, t := range d.Teachers {
						// Check that it really is a teacher, add its tid

					}
		*/
		dbdata.data.Courses = append(dbdata.data.Courses, db.Course{
			Id:       dbdata.nextId(d.Id),
			Subject:  sr,
			Groups:   glist,
			Teachers: tlist,
			//			Rooms: d.Rooms,
		})
	}
}

//TODO: Can some of that be shared with SubCourses?

// TODO: Deprecated?
func DeMultipleSubjects(w365 *W365TopLevel) {
	/* Subjects -> Subject conversion */
	// First gather keys for all Subject nodes.
	subject2key := map[W365Ref]string{}
	for _, s := range w365.Subjects {
		subject2key[s.IdStr()] = s.Shortcut
	}
	cache := map[string]W365Ref{}
	// Now check all Courses and SubCourses for multiple subjects.
	n := 0
	for i, c := range w365.Courses {
		if c.Subject == "" {
			if len(c.Subjects) == 1 {
				w365.Courses[i].Subject = c.Subjects[0]
			} else if len(c.Subjects) > 1 {
				// Make a subject name
				sklist := []string{}
				for _, sid := range c.Subjects {
					sk, ok := subject2key[sid]
					if ok {
						sklist = append(sklist, sk)
					} else {
						fmt.Printf("*ERROR* Course %s:\n  Unknown Subject: %s\n",
							c.IdStr(), sid)
					}
				}
				skname := strings.Join(sklist, ",")
				sid, ok := cache[skname]
				if !ok {
					n++
					sk := fmt.Sprintf("X%02d", n)
					sid = W365Ref(fmt.Sprintf("Id_%s", sk))
					w365.Subjects = append(w365.Subjects, Subject{
						Id:       sid,
						Name:     skname,
						Shortcut: sk,
					})
					cache[skname] = sid
					subject2key[sid] = sk

				}
				w365.Courses[i].Subject = sid
			} else if len(c.Subjects) != 0 {
				fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
					c.IdStr())
			}
		}
	}
	for i, c := range w365.SubCourses {
		if c.Subject == "" {
			if len(c.Subjects) == 1 {
				w365.SubCourses[i].Subject = c.Subjects[0]
			} else if len(c.Subjects) > 1 {
				// Make a subject name
				sklist := []string{}
				for _, sid := range c.Subjects {
					sk, ok := subject2key[sid]
					if ok {
						sklist = append(sklist, sk)
					} else {
						fmt.Printf("*ERROR* Course %s:\n  Unknown Subject: %s\n",
							c.IdStr(), sid)
					}
				}
				skname := strings.Join(sklist, ",")
				sid, ok := cache[skname]
				if !ok {
					n++
					sk := fmt.Sprintf("X%02d", n)
					sid = W365Ref(fmt.Sprintf("Id_%s", sk))
					w365.Subjects = append(w365.Subjects, Subject{
						Id:       sid,
						Name:     skname,
						Shortcut: sk,
					})
					cache[skname] = sid
					subject2key[sid] = sk

				}
				w365.SubCourses[i].Subject = sid
			} else if len(c.Subjects) != 0 {
				fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n",
					c.IdStr())
			}
		}
	}
}
