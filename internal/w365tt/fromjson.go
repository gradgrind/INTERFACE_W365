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
	return v
}

func defaultMinus1(v interface{}) int {
	if v == nil {
		return -1
	}
	return int(v.(float64))
}

type xData struct {
	w365         W365TopLevel
	data         db.DbTopLevel
	dbi          db.DbRef // counter, for db indexes
	teachers     map[W365Ref]db.DbRef
	subjects     map[W365Ref]db.DbRef
	subjectmap   map[W365Ref]string // Subject Tag (Shortcut)
	rooms        map[W365Ref]db.DbRef
	roomtag      map[db.DbRef]string    // Room Tag (Shortcut)
	roomchoices  map[string]db.DbRef    // New RoomChoiceGroup name -> db Id
	roomgroups   map[W365Ref][]db.DbRef // list of Rooms
	pregroups    map[W365Ref]string
	groups       map[W365Ref]db.DbRef
	classes      map[W365Ref]db.DbRef
	courses      map[W365Ref]db.DbRef
	subcourses   map[W365Ref]db.DbRef
	supercourses map[W365Ref]db.DbRef
	newsubjects  map[string]db.DbRef // New Subject name -> db Id
}

func LoadJSON(jsonpath string) db.DbTopLevel {
	dbdata := xData{
		w365: ReadJSON(jsonpath),
		data: db.DbTopLevel{},
		dbi:  0,
	}

	dbdata.addInfo()
	dbdata.addDays()
	dbdata.addHours()
	dbdata.addTeachers()
	dbdata.addSubjects()
	dbdata.addRooms()
	// RoomChoicesGroups: W365 has none of these – they must be generated
	// from the PreferredRooms lists of courses.
	dbdata.roomchoices = map[string]db.DbRef{}
	dbdata.data.RoomChoiceGroups = []db.RoomChoiceGroup{}
	// The RoomGroups from W365 are used by some courses. The listed rooms
	// should build the Rooms list for the course.
	dbdata.roomgroups = map[W365Ref][]db.DbRef{}
	for _, d := range dbdata.w365.RoomGroups {
		rlist := []db.DbRef{}
		for _, r := range d.Rooms {
			rr, ok := dbdata.rooms[r]
			if !ok {
				fmt.Printf("*ERROR* Unknown Room in RoomGroup %s:\n  %s\n",
					d.Id, r)
				continue
			}
			rlist = append(rlist, rr)
		}
		dbdata.roomgroups[d.Id] = rlist
	}
	dbdata.addGroups()
	dbdata.addClasses()
	dbdata.addCourses()
	dbdata.addCourses()
	dbdata.addSuperCourses()
	dbdata.addSubCourses()
	dbdata.addLessons()
	return dbdata.data
}

func (dbdata *xData) nextId() db.DbRef {
	dbdata.dbi++
	return dbdata.dbi
}

func (dbdata *xData) addInfo() {
	dbdata.data.Reference = dbdata.w365.W365TT.Scenario
	dbdata.data.Info = db.Info{
		Institution:        dbdata.w365.W365TT.SchoolName,
		Reference:          dbdata.w365.W365TT.Schedule,
		FirstAfternoonHour: dbdata.w365.W365TT.FirstAfternoonHour,
		MiddayBreak:        dbdata.w365.W365TT.MiddayBreak,
	}
}

func (dbdata *xData) addDays() {
	dbdata.data.Days = []db.Day{}
	for _, d := range dbdata.w365.Days {
		dbdata.data.Days = append(dbdata.data.Days, db.Day{
			Id:        dbdata.nextId(),
			Tag:       d.Shortcut,
			Name:      d.Name,
			Reference: d.Id,
		})
	}
}

func (dbdata *xData) addHours() {
	dbdata.data.Hours = []db.Hour{}
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
		tag := d.Shortcut
		if tag == "" {
			tag = fmt.Sprintf("(%d)", i+1)
		}
		dbdata.data.Hours = append(dbdata.data.Hours, db.Hour{
			Id:        dbdata.nextId(),
			Tag:       tag,
			Name:      d.Name,
			Start:     d.Start,
			End:       d.End,
			Reference: d.Id,
		})
	}
}

func (dbdata *xData) addTeachers() {
	dbdata.data.Teachers = []db.Teacher{}
	dbdata.teachers = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.Teachers {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		tr := dbdata.nextId()
		dbdata.data.Teachers = append(dbdata.data.Teachers, db.Teacher{
			Id:               tr,
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
			Reference:        d.Id,
		})
		dbdata.teachers[d.Id] = tr
	}
}

func (dbdata *xData) addSubjects() {
	dbdata.data.Subjects = []db.Subject{}
	dbdata.subjects = map[W365Ref]db.DbRef{}
	dbdata.newsubjects = map[string]db.DbRef{}
	dbdata.subjectmap = map[W365Ref]string{}
	for _, d := range dbdata.w365.Subjects {
		sr := dbdata.nextId()
		dbdata.data.Subjects = append(dbdata.data.Subjects, db.Subject{
			Id:        sr,
			Tag:       d.Shortcut,
			Name:      d.Name,
			Reference: d.Id,
		})
		dbdata.subjects[d.Id] = sr
		dbdata.subjectmap[d.Id] = d.Shortcut
	}
}

func (dbdata *xData) addRooms() {
	dbdata.data.Rooms = []db.Room{}
	dbdata.rooms = map[W365Ref]db.DbRef{}
	dbdata.roomtag = map[db.DbRef]string{}
	for _, d := range dbdata.w365.Rooms {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		rr := dbdata.nextId()
		dbdata.data.Rooms = append(dbdata.data.Rooms, db.Room{
			Id:           rr,
			Tag:          d.Shortcut,
			Name:         d.Name,
			NotAvailable: a,
			Reference:    d.Id,
		})
		dbdata.rooms[d.Id] = rr
		dbdata.roomtag[rr] = d.Shortcut
	}
}

func (dbdata *xData) addGroups() {
	// Every Group must be within one – and only one – Class Division.
	// To handle that, the data for the Groups is gathered here, but the
	// Elements are only added to the database when the Divisions are read.
	dbdata.data.Groups = []db.Group{}
	dbdata.pregroups = map[W365Ref]string{}
	dbdata.groups = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.Groups {
		dbdata.pregroups[d.Id] = d.Shortcut
	}
}

func (dbdata *xData) addClasses() {
	dbdata.data.Classes = []db.Class{}
	dbdata.classes = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.Classes {
		a := d.Absences
		if len(d.Absences) == 0 {
			a = []db.TimeSlot{}
		}
		// Get the divisions and add their groups to the database.
		divs := []db.Division{}
		for _, wdiv := range d.Divisions {
			glist := []db.DbRef{}
			for _, g := range wdiv.Groups {
				gtag, ok := dbdata.pregroups[g] // get Tag
				if ok {
					// Add Group to database
					if _, nok := dbdata.groups[g]; nok {
						fmt.Printf("*ERROR* Group Defined in"+
							" multiple Divisions:\n  -- %s\n", g)
					}
					gr := dbdata.nextId()
					dbdata.data.Groups = append(dbdata.data.Groups, db.Group{
						Id:        gr,
						Tag:       gtag,
						Reference: g,
					})
					dbdata.groups[g] = gr
					glist = append(glist, gr)
				} else {
					fmt.Printf("*ERROR* Unknown Group in Class %s,"+
						" Division %s:\n  %s\n", d.Shortcut, wdiv.Name, g)
				}
			}
			divs = append(divs, db.Division{
				Name:      wdiv.Name,
				Groups:    glist,
				Reference: wdiv.Id,
			})
		}
		cr := dbdata.nextId()
		dbdata.data.Classes = append(dbdata.data.Classes, db.Class{
			Id:               cr,
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
			Reference:        d.Id,
		})
		dbdata.classes[d.Id] = cr
	}
}

func (dbdata *xData) readCourse(
	id W365Ref,
	subject W365Ref,
	subjects []W365Ref,
	groups []W365Ref,
	teachers []W365Ref,
	rooms []W365Ref,
) (db.DbRef, []db.DbRef, []db.DbRef, []db.DbRef) {
	// Deal with subject
	var sr db.DbRef = 0
	var ok bool
	msg := "*ERROR* Course %s:\n  Unknown Subject: %s\n"
	if subject == "" {
		if len(subjects) == 1 {
			wsid := subjects[0]
			sr, ok = dbdata.subjects[wsid]
			if !ok {
				fmt.Printf(msg, id, wsid)
			}
		} else if len(subjects) > 1 {
			// Make a subject name
			sklist := []string{}
			for _, wsid := range subjects {
				// Need Shortcut field
				sk, ok := dbdata.subjectmap[wsid]
				if ok {
					sklist = append(sklist, sk)
				} else {
					fmt.Printf(msg, id, wsid)
				}
			}
			skname := strings.Join(sklist, ",")
			sr, ok = dbdata.newsubjects[skname]
			if !ok {
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
			fmt.Printf("*ERROR* Course has both Subject AND Subjects: %s\n", id)
		}
		wsid := subject
		sr, ok = dbdata.subjects[wsid]
		if !ok {
			fmt.Printf(msg, id, wsid)
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
				fmt.Printf("*ERROR* Unknown group in Course %s:\n  %s\n", id, g)
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
			fmt.Printf("*ERROR* Unknown teacher in Course %s:\n  %s\n", id, t)
			continue
		}
		tlist = append(tlist, tr)
	}
	// Deal with rooms. W365 can have a single RoomGroup or a list of Rooms
	rclist := []db.DbRef{} // choice list
	var rlist []db.DbRef   // actual room list
	for _, r := range rooms {
		rr, ok := dbdata.rooms[r]
		if ok {
			rclist = append(rclist, rr)
		} else {
			rl, ok := dbdata.roomgroups[r]
			if ok {
				rlist = rl
				if len(rooms) != 1 {
					rclist = []db.DbRef{}
					fmt.Printf(
						"*ERROR* Mixed Rooms and RoomGroups in Course %s\n", id)
				}
				break
			} else {
				fmt.Printf("*ERROR* Unknown room in Course %s:\n  %s\n", id, r)
				continue
			}
		}
	}
	if len(rclist) == 1 {
		// Make the Room compulsory.
		rlist = rclist
	} else if len(rclist) > 1 {
		// Need a RoomChoiceGroup.
		// Reuse these if the same list appears again, but treat the
		// order as significant.
		rslist := []string{}
		for _, r := range rclist {
			rslist = append(rslist, dbdata.roomtag[r])
		}
		rs := strings.Join(rslist, ",")
		rr, ok := dbdata.roomchoices[rs]
		if !ok {
			rk := fmt.Sprintf("RC%03d", len(dbdata.roomchoices)+1)
			rr = dbdata.nextId()
			dbdata.data.RoomChoiceGroups = append(
				dbdata.data.RoomChoiceGroups, db.RoomChoiceGroup{
					Id:    rr,
					Tag:   rk,
					Name:  rs,
					Rooms: rclist,
				})
			dbdata.roomchoices[rs] = rr
		}
		rlist = []db.DbRef{rr}
	}
	return sr, glist, tlist, rlist
}

func (dbdata *xData) addCourses() {
	dbdata.data.Courses = []db.Course{}
	dbdata.courses = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.Courses {
		sr, glist, tlist, rlist := dbdata.readCourse(
			d.Id, d.Subject, d.Subjects, d.Groups, d.Teachers, d.PreferredRooms)
		cr := dbdata.nextId()
		dbdata.data.Courses = append(dbdata.data.Courses, db.Course{
			Id:        cr,
			Subject:   sr,
			Groups:    glist,
			Teachers:  tlist,
			Rooms:     rlist,
			Reference: d.Id,
		})
		dbdata.courses[d.Id] = cr
	}
}

func (dbdata *xData) addSuperCourses() {
	dbdata.data.SuperCourses = []db.SuperCourse{}
	dbdata.supercourses = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.SuperCourses {
		cr := dbdata.nextId()
		sr, ok := dbdata.subjects[d.Subject]
		if !ok {
			fmt.Printf("*ERROR* Unknown Subject in SuperCourse %s:\n  %s\n",
				d.Id, d.Subject)
			continue
		}
		dbdata.data.SuperCourses = append(dbdata.data.SuperCourses, db.SuperCourse{
			Id:        cr,
			Subject:   sr,
			Reference: d.Id,
		})
		dbdata.supercourses[d.Id] = cr
	}
}

func (dbdata *xData) addSubCourses() {
	dbdata.data.SubCourses = []db.SubCourse{}
	dbdata.subcourses = map[W365Ref]db.DbRef{}
	for _, d := range dbdata.w365.SubCourses {
		sr, glist, tlist, rlist := dbdata.readCourse(
			d.Id, d.Subject, d.Subjects, d.Groups, d.Teachers, d.PreferredRooms)
		sc, ok := dbdata.supercourses[d.SuperCourse]
		if !ok {
			fmt.Printf("*ERROR* Unknown SuperCourse in SubCourse %s:\n  %s\n",
				d.Id, d.SuperCourse)
			continue
		}
		cr := dbdata.nextId()
		dbdata.data.SubCourses = append(dbdata.data.SubCourses, db.SubCourse{
			Id:          cr,
			SuperCourse: sc,
			Subject:     sr,
			Groups:      glist,
			Teachers:    tlist,
			Rooms:       rlist,
			Reference:   d.Id,
		})
		dbdata.subcourses[d.Id] = cr
	}
}

func (dbdata *xData) addLessons() {
	dbdata.data.Lessons = []db.Lesson{}
	for _, d := range dbdata.w365.Lessons {
		// The course can be either a Course or a SubCourse.
		crs, ok := dbdata.courses[d.Course]
		if !ok {
			crs, ok = dbdata.subcourses[d.Course]
			if !ok {
				fmt.Printf("*ERROR* Invalid course in Lesson %s:\n  -- %s\n",
					d.Id, d.Course)
				continue
			}
		}
		rlist := []db.DbRef{}
		for _, r := range d.LocalRooms {
			rr, ok := dbdata.rooms[r]
			if ok {
				rlist = append(rlist, rr)
			} else {
				fmt.Printf("*ERROR* Invalid room in Lesson %s:\n  -- %s\n",
					d.Id, r)
			}
		}
		dbdata.data.Lessons = append(dbdata.data.Lessons, db.Lesson{
			Id:        dbdata.nextId(),
			Course:    crs,
			Duration:  d.Duration,
			Day:       d.Day,
			Hour:      d.Hour,
			Fixed:     d.Fixed,
			Rooms:     rlist,
			Reference: d.Id,
		})
	}
}
